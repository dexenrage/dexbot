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

import "fmt"

const (
	sqlCreateTable = `
	CREATE TABLE IF NOT EXISTS users
	(
		id serial PRIMARY KEY,
		user_id bigint NOT NULL,
		item text,
		price text
	);`

	sqlGetPrices = `SELECT user_id, item, price FROM users ORDER BY user_id, id;`
)

func sqlInsertItem(userID, item, price string) string {
	const insert = `INSERT INTO users(user_id, item, price) VALUES(%s);`
	values := fmt.Sprintf(`%s,'%s','%s'`, userID, item, price)
	return fmt.Sprintf(insert, values)
}

func sqlGetItemList(userID string) string {
	const selectQuery = `SELECT item, price FROM users WHERE user_id = %s ORDER BY id;`
	return fmt.Sprintf(selectQuery, userID)
}

func sqlSetPrice(price, item, userID string) string {
	const update = `UPDATE users SET price = '%s' WHERE item = '%s' AND user_id = '%s';`
	return fmt.Sprintf(update, price, item, userID)
}

func sqlDeleteItem(userID, item string) string {
	const delete = `DELETE FROM users WHERE user_id = %s AND item = '%s';`
	return fmt.Sprintf(delete, userID, item)
}
