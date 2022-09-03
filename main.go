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
	"dexbot/catcherr"
	"dexbot/commands"
	"dexbot/config"
	"dexbot/tracker"
	"time"

	tb "gopkg.in/telebot.v3"
)

func main() {
	defer catcherr.Recover(`main`)

	settings := tb.Settings{
		Token:     config.String(`bot_token`),
		Poller:    &tb.LongPoller{Timeout: 15 * time.Second},
		ParseMode: tb.ModeMarkdown,
	}
	bot, err := tb.NewBot(settings)
	catcherr.HandleError(err)

	go tracker.Start(bot)
	commands.Handle(bot)
	bot.Start()
}
