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

package database

import (
	"context"
	"database/sql"
	"dexbot/catcherr"
	"dexbot/config"
	"fmt"
	"net/url"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var db *bun.DB

func init() {
	ctx := context.Background()

	var (
		username = config.String(`db_user`)
		password = config.String(`db_pass`)
		sslmode  = config.String(`db_ssl`)
		path     = config.String(`db_name`)

		user = url.UserPassword(username, password)
	)

	var host string
	{
		h := config.String(`db_host`)
		p := config.String(`db_port`)
		host = fmt.Sprint(h, `:`, p)
	}

	dsn := url.URL{
		Scheme: "postgres",
		Host:   host,
		User:   user,
		Path:   path,
	}
	{
		q := dsn.Query()
		q.Add(`sslmode`, sslmode)
		dsn.RawQuery = q.Encode()
	}

	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn.String())))

	// Create a Bun db on top of it.
	db = bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Create items table if not exists
	_, err := db.NewCreateTable().Model((*Item)(nil)).IfNotExists().Exec(ctx)
	catcherr.HandleError(err)
}

func AddItem(ctx context.Context, userID int64, path string, price float64) error {
	item := &Item{UserID: userID, ItemURL: path, Price: price}
	_, err := db.NewInsert().Model(item).Exec(ctx)
	return err
}

func GetItemList(ctx context.Context, userID int64) (list []Item, err error) {
	q := db.NewSelect().Model(&list).Where(`id = ?`, userID)
	err = q.Column(`item_url`, `price`).Order(`i.created_at ASC`).Scan(ctx)
	return list, err
}

func GetAllItems(ctx context.Context) (list []Item, err error) {
	err = db.NewSelect().Model(&list).Scan(ctx)
	return list, err
}

func UpdatePrice(ctx context.Context, userID int64, item string, price float64) error {
	i := Item{UserID: userID, ItemURL: item, Price: price}
	_, err := db.NewUpdate().Model(&i).WherePK().Exec(ctx)
	return err
}
func DeleteItem(ctx context.Context, userID int64, item string) (err error) {
	i := Item{UserID: userID, ItemURL: item}
	q := db.NewDelete().Model(&i).Where(`id = ?`, userID).Where(`item_url = ?`, item)
	_, err = q.Exec(ctx)
	return err
}
