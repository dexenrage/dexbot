/*
   Copyright 2022 dexenrage

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// The package main initializes Telegram bot
package main

import (
	"errors"
	"fmt"
	"log"
	u "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

var (
	g_logs *zap.SugaredLogger
	g_conf *koanf.Koanf
)

func init() {
	err := initLogger()
	if err != nil {
		log.Fatalln(err)
	}

	err = configure()
	if err != nil {
		g_logs.Fatal(err)
	}

	err = checkDB()
	if err != nil {
		g_logs.Fatal(err)
	}
}

func initLogger() error {
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}
	logsFolder := filepath.Join(workdir, `logs`)

	err = os.Mkdir(logsFolder, 0766)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		g_logs.Fatalf(errCreateDir, logsFolder, err)
	}

	logsFile := fmt.Sprint(`log-`, time.Now().Format(time.RFC3339))
	logsFilePath := filepath.Join(logsFolder, logsFile)

	logConf := zap.NewDevelopmentConfig()
	logConf.OutputPaths = []string{logsFilePath}

	logs, err := logConf.Build()
	if err != nil {
		errText := fmt.Sprintf(errZapInit, err)
		err = errors.New(errText)
		return err
	}
	g_logs = logs.Sugar()
	return err
}

func configure() error {
	defer g_logs.Sync()

	cfgFile, err := findInParentDirs(3, `config`, `config.yml`)
	if err != nil {
		return err
	}
	fp := file.Provider(cfgFile)

	g_conf = koanf.New(`.`)
	if err := g_conf.Load(fp, yaml.Parser()); err != nil {
		g_logs.Fatalf(errLoadingConfig, err)
	}
	return err
}

func cssElementsToSlice() []string {
	raw := g_conf.String(`css_elements`)
	return strings.Split(raw, ` `)
}

func allowedLinksToSlice() []string {
	raw := g_conf.String(`allowed_links`)
	return strings.Split(raw, ` `)
}

func findInParentDirs(maxParentDirs int, items ...string) (path string, err error) {
	defer g_logs.Sync()

	workdir, err := os.Getwd()
	if err != nil {
		g_logs.Fatalf(errGetWorkDir, err)
	}

	var itemPath string
	for i := range items {
		itemPath = filepath.Join(itemPath, items[i])
	}
	path = filepath.Join(workdir, itemPath)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		for i := 0; i < maxParentDirs; i++ {
			workdir = filepath.Dir(workdir) // Go to parent directory
			path = filepath.Join(workdir, itemPath)

			_, err = os.Stat(path)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return ``, err
	}
	return path, err
}

func main() {
	defer g_logs.Sync()

	bot, err := tb.NewBot(tb.Settings{
		Token:     g_conf.String(`bot_token`),
		Poller:    &tb.LongPoller{Timeout: 15 * time.Second},
		ParseMode: tb.ModeMarkdown,
	})
	if err != nil {
		g_logs.Fatal(err)
	}

	go priceNotifier(*bot)
	handleCMD(*bot)
	bot.Start()
}

func handleCMD(bot tb.Bot) {
	const (
		startCMD  = `/start`
		helpCMD   = `/help`
		addCMD    = `/add`
		listCMD   = `/list`
		deleteCMD = `/rm`
	)

	bot.Handle(startCMD, helpMessage)
	bot.Handle(helpCMD, helpMessage)
	bot.Handle(addCMD, addItem)
	bot.Handle(listCMD, sendList)
	bot.Handle(deleteCMD, deleteItem)
}

func helpMessage(msg tb.Context) error {
	return msg.Send(helpMSG())
}

func findInAllowedURLs(url string) (string, error) {
	var err error
	links := allowedLinksToSlice()
	for _, v := range links {
		if strings.HasPrefix(url, v) {
			return url, err
		}
	}
	errText := fmt.Sprintf(errURLNotAllowed, url)
	err = errors.New(errText)
	return url, err
}

/* Trims the URL protocol to save space.
 * For example, for a more beautiful message to the user,
 * because it is optional to specify the URL protocol in Telegram.*/
func trimProtocol() string {
	url := allowedLinksToSlice()[0]
	url = strings.TrimPrefix(url, `https://`)
	url = strings.TrimPrefix(url, `http://`)
	return url
}

func getURLPath(url string) (string, error) {
	u, err := u.Parse(url)
	if err != nil {
		return u.Path, err
	}
	return u.Path, err
}

func addItem(msg tb.Context) error {
	defer g_logs.Sync()

	if len(msg.Args()) == 0 {
		return msg.Send(addErrMSG)
	}

	url := msg.Args()[0]
	if !valid.IsURL(url) {
		return msg.Send(addErrMSG)
	}

	url, err := findInAllowedURLs(url)
	if err != nil {
		g_logs.Error(err)
		return msg.Send(addErrMSG)
	}

	price, err := parse(url)
	if err != nil {
		g_logs.Error(err)
		return msg.Send(addErrMSG)
	}

	item, err := getURLPath(url)
	if err != nil {
		g_logs.Error(err)
		return msg.Send(addErrMSG)
	}

	err = addToDB(msg.Sender().ID, item, price)
	if err != nil {
		g_logs.Error(err)
		return msg.Send(intErrMSG)
	}
	return msg.Send(addedMSG)
}

func sendList(msg tb.Context) error {
	defer g_logs.Sync()

	item, _, err := getItems(msg.Sender().ID)
	if err != nil {
		g_logs.Error(err)
		return msg.Send(intErrMSG)
	}

	if len(item) == 0 {
		return msg.Send(listEmptyMSG)
	}

	list := listHeaderMSG
	{
		var (
			prefix = trimProtocol()
			itemID int
		)
		for i, url := range item {
			itemID = i + 1
			list += fmt.Sprint(itemID, ". ", prefix, url, "\n")
		}
	}
	return msg.Send(list, tb.NoPreview)
}

type priceStruct struct {
	UserID64 int64
	Message  string
}

func priceNotifier(msg tb.Bot) {
	priceChan := make(chan priceStruct)
	defer close(priceChan)

	go tracker(priceChan)

	var user tb.User
	for out := range priceChan {
		user.ID = out.UserID64
		msg.Send(&user, out.Message)
	}
}

func deleteItem(msg tb.Context) error {
	defer g_logs.Sync()

	if len(msg.Args()) == 0 {
		return msg.Send(rmErrMSG)
	}
	itemNumber := msg.Args()[0]

	// We receive a string type from the user and check that the string
	// contains only a number and whether it contains it at all.
	if !valid.IsInt(itemNumber) {
		return msg.Send(rmErrMSG)
	}

	num, err := strconv.Atoi(itemNumber) // Initially, we get a string type and need to convert it.
	if num <= 0 {
		errText := fmt.Sprintf(errOutOfRange, err)
		err = errors.New(errText) // Prevents runtime error.
	}
	if err != nil {
		g_logs.Error(err)
		return msg.Send(rmErrMSG)
	}

	list, _, err := getItems(msg.Sender().ID)
	if num > len(list) {
		errText := fmt.Sprintf(errOutOfRange, err)
		err = errors.New(errText) // Prevents runtime error.
	}
	if err != nil {
		g_logs.Error(err)
		return msg.Send(rmErrMSG)
	}

	// Numbering starts at 0, but the user gets a list in which numbering starts at 1.
	item := list[num-1]

	err = deleteFromDB(msg.Sender().ID, item)
	if err != nil {
		g_logs.Error(err)
		return msg.Send(intErrMSG)
	}
	return msg.Send(removedMSG)
}
