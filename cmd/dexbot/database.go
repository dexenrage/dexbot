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
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func connect() *pgx.Conn {
	defer g_logs.Sync()
	const dbConfigFields = `user=%s password=%s host=%s port=%s dbname=%s`

	var (
		user   = g_conf.String(`db_user`)
		pass   = g_conf.String(`db_pass`)
		host   = g_conf.String(`db_host`)
		port   = g_conf.String(`db_port`)
		dbName = g_conf.String(`db_name`)

		dbConfigSource = fmt.Sprintf(dbConfigFields, user, pass, host, port, dbName)
	)

	dbConnConfig, err := pgx.ParseConfig(dbConfigSource)
	if err != nil {
		g_logs.Fatalf(errParsingDatabaseURI, err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), dbConnConfig)
	if err != nil {
		g_logs.Fatal(errDatabaseConnection, err)
	}
	return conn
}

func checkDB() error {
	conn := connect()
	defer conn.Close(context.Background())

	resp, err := conn.Query(context.Background(), sqlCreateTable)
	if err != nil {
		return err
	}
	defer resp.Close()
	return err
}

func addToDB(userID64 int64, item, price string) error {
	conn := connect()
	defer conn.Close(context.Background())

	userID := fmt.Sprint(userID64)
	query := sqlInsertItem(userID, item, price)

	resp, err := conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer resp.Close()
	return err
}

func getItems(userID64 int64) (item, price []string, err error) {
	conn := connect()
	defer conn.Close(context.Background())

	var (
		userID = fmt.Sprint(userID64)
		query  = sqlGetItemList(userID)
	)

	resp, err := conn.Query(context.Background(), query)
	if err != nil {
		return item, price, err
	}
	defer resp.Close()

	for resp.Next() {
		v, err := resp.Values()
		if err != nil {
			return item, price, err
		}
		item = append(item, v[0].(string))
		price = append(price, v[1].(string))
	}
	return item, price, err
}

func getPrices() (userID64 []int64, item, price []string, err error) {
	conn := connect()
	defer conn.Close(context.Background())

	resp, err := conn.Query(context.Background(), sqlGetPrices)
	if err != nil {
		return userID64, item, price, err
	}
	defer resp.Close()

	for resp.Next() {
		v, err := resp.Values()
		if err != nil {
			return userID64, item, price, err
		}
		userID64 = append(userID64, v[0].(int64))
		item = append(item, v[1].(string))
		price = append(price, v[2].(string))
	}
	return userID64, item, price, err
}

func updatePrice(userID64 int64, item, price string) error {
	conn := connect()
	defer conn.Close(context.Background())

	userID := fmt.Sprint(userID64)
	query := sqlSetPrice(price, item, userID)

	resp, err := conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer resp.Close()
	return err
}

func deleteFromDB(userID64 int64, item string) error {
	conn := connect()
	defer conn.Close(context.Background())

	userID := fmt.Sprint(userID64)
	query := sqlDeleteItem(userID, item)

	resp, err := conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer resp.Close()
	return err
}
