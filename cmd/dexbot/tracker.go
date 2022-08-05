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

package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func tracker(priceChan chan<- priceStruct) {
	defer g_logs.Sync()

	userID, item, price, err := getPrices()
	if err != nil {
		g_logs.Error(err)
	}

	var (
		urlPrefix = trimProtocol()

		currentUser int64
		idInt       int
		items       []string

		parsedPrice      string
		parsedPriceFloat float64
		savedPriceFloat  float64

		itemID   string
		itemURL  string
		oldPrice string
		newPrice string
		message  string

		currency = g_conf.String(`currency`)
	)

	// Strings to format by fmt.Sprintf()
	// The `*` character is needed to make text bold
	// in Telegram Markdown style. For example: *bold text*
	const (
		itemID_F   = `*%s%d*`    // idMSG, idInt
		itemURL_F  = `%s*%s%s*`  // urlMSG, urlPrefix, url
		oldPrice_F = `%s%s %s`   // oldPriceMSG, price[i], currency
		newPrice_F = `*%s%s %s*` // newPriceMSG, parsedPrice, currency
	)

	for i, url := range item {
		parsedPrice, err = parse(url)
		if err != nil {
			g_logs.Error(err)
			continue
		}

		// In Go float, numbers are separated by a dot, not a comma.
		parsedPrice = strings.ReplaceAll(parsedPrice, `,`, `.`)
		parsedPriceFloat, err = strconv.ParseFloat(parsedPrice, 64)
		if err != nil {
			g_logs.Error(err)
			continue
		}

		savedPriceFloat, err = strconv.ParseFloat(price[i], 64)
		if err != nil {
			g_logs.Error(err)
			continue
		}

		/*	A small check so as not to create unnecessary queries
			to the database if the user ID has not changed */
		if currentUser != userID[i] {
			items, _, err = getItems(userID[i])
			if err != nil {
				g_logs.Error(err)
				continue
			}
		}

		if savedPriceFloat != parsedPriceFloat {
			err = updatePrice(userID[i], url, parsedPrice)
			if err != nil {
				g_logs.Error(err)
				continue
			}

			if savedPriceFloat < parsedPriceFloat {
				message = priceUpMSG

			}
			if savedPriceFloat > parsedPriceFloat {
				message = priceDownMSG
			}

			// To add an item ID from a user's list to a message
			for index, value := range items {
				if value == url {
					idInt = index + 1
					break
				}
			}

			itemID = fmt.Sprintf(itemID_F, idMSG, idInt)
			itemURL = fmt.Sprintf(itemURL_F, urlMSG, urlPrefix, url)
			oldPrice = fmt.Sprintf(oldPrice_F, oldPriceMSG, price[i], currency)
			newPrice = fmt.Sprintf(newPrice_F, newPriceMSG, parsedPrice, currency)

			message = fmt.Sprint(message, itemID, itemURL, oldPrice, newPrice)
			priceChan <- priceStruct{userID[i], message}
		}
	}

	timer() // Sleep until next check.
	go tracker(priceChan)
}

func timer() {
	d := g_conf.String(`duration`)

	timeDuration, err := time.ParseDuration(d)
	if err != nil {
		g_logs.Error(err)
	}
	duration := time.Now().Add(timeDuration)

	time.Sleep(time.Until(duration))
}

func parse(url string) (price string, err error) {
	url, err = getURLPath(url)
	if err != nil {
		return price, err
	}
	url = fmt.Sprint(allowedLinksToSlice()[0], url)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return price, err
	}

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return price, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return price, err
	}
	if resp.StatusCode != 200 {
		errText := fmt.Sprintf(errHTTPStatusCode, url, resp.StatusCode)
		err = errors.New(errText)
		return price, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return price, err
	}

	for _, element := range cssElementsToSlice() {
		doc.Find(element).Each(func(num int, s *goquery.Selection) {
			price = s.Text()
		})
		if price != `` {
			break
		}
	}

	const floatNumbers = `[+-]?([0-9]*[,])?[0-9]+`
	r, err := regexp.Compile(floatNumbers)
	if err != nil {
		return price, err

	}
	price = r.FindString(price)

	// If an empty string was parsed (for example, when the css element does not exist)
	if price == `` {
		errText := fmt.Sprintf(errCannotParsePrice, url)
		err = errors.New(errText)
		return price, err
	}

	/*	In Go float, numbers are separated by a dot, not a comma,
		so we replace the sign so that later it will be easier
		to convert string to float */
	price = strings.ReplaceAll(price, `,`, `.`)
	return price, err
}
