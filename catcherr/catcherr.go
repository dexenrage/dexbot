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

// The package catcherr handles errors
package catcherr

import (
	"errors"
	"fmt"
	"log"

	tb "gopkg.in/telebot.v3"
)

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleErrorChannel(errChan chan<- error, err error) {
	if err != nil {
		errChan <- err
	}
}

func HandleAndResponse(msg tb.Context, message string, err error) {
	if err != nil {
		LogError(`catcherr.HandleAndResponse()`, msg.Send(message))
		HandleError(err)
	}
}

func LogError(sender string, err error) {
	if err != nil {
		const tmpl = `[ Sender: %s ]: %v `
		log.Printf(tmpl, sender, err)
	}
}

func Recover(sender string) {
	if r := recover(); r != nil {
		const tmpl = `[ Sender: %s ]: %v `
		log.Printf(tmpl, sender, r)
	}
}

func RecoverAndReturnError() (err error) {
	if r := recover(); r != nil {
		return errors.New(fmt.Sprint(r))
	}
	return nil
}
