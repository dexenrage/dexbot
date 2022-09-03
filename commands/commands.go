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

package commands

import (
	"context"
	"dexbot/actions"
	"dexbot/catcherr"
	"dexbot/database"
	"dexbot/messages"
	"fmt"
	"net/url"
	"strconv"
	"time"

	tb "gopkg.in/telebot.v3"
)

func Handle(bot *tb.Bot) {
	const (
		startCMD  = `/start`
		helpCMD   = `/help`
		addCMD    = `/add`
		listCMD   = `/list`
		deleteCMD = `/rm`
	)

	bot.Handle(startCMD, help)
	bot.Handle(helpCMD, help)
	bot.Handle(addCMD, add)
	bot.Handle(listCMD, list)
	bot.Handle(deleteCMD, delete)
}

func defaultContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

func help(msg tb.Context) error { return msg.Send(messages.Help) }

func add(msg tb.Context) error {
	defer catcherr.Recover(`commands.add`)

	ctx, cancel := defaultContextTimeout()
	defer cancel()

	u, err := url.ParseRequestURI(msg.Args()[0])
	catcherr.HandleAndResponse(msg, messages.NeedCorrectLink, err)
	path := u.String()

	if !isAllowedURL(path) {
		return msg.Send(messages.NeedCorrectLink)
	}

	price, err := actions.GetPrice(ctx, path)
	catcherr.HandleAndResponse(msg, messages.NeedCorrectLink, err)

	err = database.AddItem(ctx, msg.Sender().ID, path, price)
	catcherr.HandleAndResponse(msg, messages.InternalError, err)

	return msg.Send(messages.AddedSuccessfully)
}

func list(msg tb.Context) error {
	defer catcherr.Recover(`commands.list`)

	ctx, cancel := defaultContextTimeout()
	defer cancel()

	list, err := database.GetItemList(ctx, msg.Sender().ID)
	catcherr.HandleAndResponse(msg, messages.InternalError, err)
	if len(list) == 0 {
		return msg.Send(messages.EmptyList)
	}

	message := messages.ListHeader
	for i, v := range list {
		link := actions.TrimURLScheme(v.ItemURL)
		message += fmt.Sprint(i+1, ". ", link, "\n")
	}
	return msg.Send(message, tb.NoPreview)
}

func delete(msg tb.Context) error {
	defer catcherr.Recover(`commands.delete`)

	ctx, cancel := defaultContextTimeout()
	defer cancel()

	num, err := strconv.Atoi(msg.Args()[0])
	catcherr.HandleAndResponse(msg, messages.RemoveError, err)

	if num <= 0 {
		return msg.Send(messages.RemoveError)
	}

	list, err := database.GetItemList(ctx, msg.Sender().ID)
	catcherr.HandleAndResponse(msg, messages.RemoveError, err)
	if num > len(list) {
		return msg.Send(messages.RemoveError)
	}

	// Numbering starts at 0, but the user gets a list in which numbering starts at 1.
	item := (list)[num-1].ItemURL

	err = database.DeleteItem(ctx, msg.Sender().ID, item)
	catcherr.HandleAndResponse(msg, messages.RemoveError, err)

	return msg.Send(messages.Removed)
}
