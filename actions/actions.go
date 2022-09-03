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

package actions

import (
	"context"
	"dexbot/catcherr"
	"dexbot/config"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TrimURLScheme(path string) string {
	u, err := url.Parse(path)
	catcherr.HandleError(err)

	prefix := fmt.Sprint(u.Scheme, `://`)
	return strings.TrimPrefix(path, prefix)
}

func GetPrice(ctx context.Context, path string) (price float64, err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	jar, err := cookiejar.New(nil)
	catcherr.HandleError(err)

	client := http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	catcherr.HandleError(err)

	resp, err := client.Do(req)
	catcherr.HandleError(err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := catcherr.HTTPStatusCode(http.StatusOK, resp.StatusCode)
		catcherr.HandleError(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	catcherr.HandleError(err)

	elems := config.StringSlice(`css_elements`, config.DefaultSeparator)

	var priceString string
	for i := range elems {
		doc.Find(elems[i]).Each(func(i int, s *goquery.Selection) {
			priceString = s.Text()
		})

		if len(priceString) != 0 {
			break
		}
	}

	const floatNumbers = `([0-9]*[,]|[.])?[0-9]+`
	r, err := regexp.Compile(floatNumbers)
	catcherr.HandleError(err)

	priceString = r.FindString(priceString)

	if len(priceString) == 0 {
		catcherr.HandleError(catcherr.MissingCSSElement())
	}

	/*	In Go float, numbers are separated by a dot, not a comma,
		so we replace the sign so that later it will be easier
		to convert string to float */
	priceString = strings.ReplaceAll(priceString, `,`, `.`)

	price, err = strconv.ParseFloat(priceString, 64)
	catcherr.HandleError(err)
	return price, err
}
