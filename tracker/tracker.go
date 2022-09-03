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

package tracker

import (
	"context"
	"dexbot/actions"
	"dexbot/catcherr"
	"dexbot/config"
	"dexbot/database"
	"dexbot/messages"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	tb "gopkg.in/telebot.v3"
)

type priceData struct {
	UserID       int64
	ItemURL      string
	OldPrice     float64
	CurrentPrice float64
}

func Start(bot *tb.Bot) {
	const errorSender = `tracker.Start()`
	defer catcherr.Recover(errorSender)

	var (
		duration = config.String(`duration`)
		ctx      = context.Background()
	)

	for {
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error { return timer(duration) })

		g.Go(func() error {
			data := tracker(ctx)
			for _, v := range data {
				if v.OldPrice == v.CurrentPrice {
					continue
				}

				itemList, err := database.GetItemList(ctx, v.UserID)
				if err != nil {
					catcherr.LogError(errorSender, err)
					continue
				}

				err = database.UpdatePrice(ctx, v.UserID, v.ItemURL, v.CurrentPrice)
				if err != nil {
					catcherr.LogError(errorSender, err)
					continue
				}

				msg := prepareMessage(v.ItemURL, v.OldPrice, v.CurrentPrice, itemList)
				_, err = bot.Send(&tb.User{ID: v.UserID}, msg, tb.NoPreview)
				catcherr.LogError(errorSender, err)
			}
			return nil
		})
		catcherr.HandleError(g.Wait())
	}
}

func prepareMessage(

	itemURL string,
	oldPrice float64,
	currentPrice float64,
	itemList []database.Item,

) (message string) {

	var itemID int
	for i, v := range itemList {
		if v.ItemURL == itemURL {
			itemID = i + 1
			break
		}
	}

	var priceStatus string
	switch {
	case oldPrice < currentPrice:
		priceStatus = messages.PriceUp
	case oldPrice > currentPrice:
		priceStatus = messages.PriceDown
	}

	message = fmt.Sprintf(
		messages.ChangedPriceTemplate,
		priceStatus,
		itemID,
		actions.TrimURLScheme(itemURL),
		oldPrice,
		currentPrice,
	)
	return message
}

func tracker(ctx context.Context) (data []priceData) {
	items, err := database.GetAllItems(ctx)
	catcherr.HandleError(err)

	for _, v := range items {
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() (err error) {
			defer func() { err = catcherr.RecoverAndReturnError() }()

			price, err := actions.GetPrice(ctx, v.ItemURL)
			catcherr.HandleError(err)

			data = append(data, priceData{
				UserID:       v.UserID,
				ItemURL:      v.ItemURL,
				OldPrice:     v.Price,
				CurrentPrice: price,
			})
			return err
		})
		catcherr.LogError(`tracker.tracker()`, g.Wait())
		time.Sleep(1 * time.Second) // To avoid HTTP request flood
	}
	return data
}

func timer(duration string) (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	d, err := time.ParseDuration(duration)
	catcherr.HandleError(err)

	dur := time.Now().Add(d)
	time.Sleep(time.Until(dur))
	return err
}
