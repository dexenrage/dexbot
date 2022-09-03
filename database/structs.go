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
	"time"

	"github.com/uptrace/bun"
)

type Item struct {
	bun.BaseModel `bun:"table:items,alias:i"`
	UserID        int64  `bun:"id,notnull"`
	ItemURL       string `bun:",notnull"`
	Price         float64
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
