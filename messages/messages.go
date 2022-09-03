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

package messages

import (
	"dexbot/config"
	"fmt"
)

const (
	ChangedPriceTemplate = `
	%s
	📍 ID: *%d*
	🔗 *%s*
	
	▫ Старая: %f
	🔥 Новая: *%f*`

	PriceUp   = "❌ Цена выросла"
	PriceDown = "✅ Цена упала"
)

const (
	AddedSuccessfully = `✅ Товар успешно добавлен в трекер.`
	NeedCorrectLink   = "❌ Пожалуйста, отправьте правильную ссылку на товар.\n🔗 Используйте */add <url>*"

	ListHeader = "📝 Список отслеживаемых товаров:\n"
	EmptyList  = "📝 Список пуст.\n🔗 Используйте */add <url>* чтобы добавить товары."

	Removed     = "✅ Товар успешно удалён из трекера."
	RemoveError = `❌ Пожалуйста, отправьте правильный ID товара.
🔗 Используйте */rm <id>*

📝 Если Вы не знаете нужный ID - введите */list*.`

	InternalError = "❌ Произошла внутренняя ошибка.\n⏳ Ожидайте, скоро всё заработает."
)

func Help() string {
	name := config.String(`bot_name`)

	helpMSG := `👤 *%s* 👤

/help - Показать это сообщение.
/add - Добавить в трекер.
/list - Список товаров.
/rm - Удалить из трекера.

🔰 Выгодных покупок! 🔰`
	return fmt.Sprintf(helpMSG, name)
}
