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

func helpMSG() string {
	name := g_conf.String(`bot_name`)

	helpMSG := `👤 *` + name + `* 👤

/help - Показать это сообщение.
/add - Добавить в трекер.
/list - Список товаров.
/rm - Удалить из трекера.

🔰 Выгодных покупок! 🔰`
	return helpMSG
}

const (
	addedMSG  = "✅ Товар успешно добавлен в трекер."
	addErrMSG = "❌ Пожалуйста, отправьте правильную ссылку на товар.\n🔗 Используйте */add <url>*"

	listHeaderMSG = "📝 Список отслеживаемых товаров:\n"
	listEmptyMSG  = "📝 Список пуст.\n🔗 Используйте */add <url>* чтобы добавить товары."

	removedMSG = "✅ Товар успешно удалён из трекера."
	rmErrMSG   = `❌ Пожалуйста, отправьте правильный ID товара.
🔗 Используйте */rm <id>*

📝 Если Вы не знаете нужный ID - введите */list*.`

	priceUpMSG   = "❌ Цена выросла\n"
	priceDownMSG = "✅ Цена упала\n"
	idMSG        = "\n📍 ID: "
	urlMSG       = "\n🔗 "
	oldPriceMSG  = "\n\n▫ Старая: "
	newPriceMSG  = "\n🔥 Новая: "

	intErrMSG = "❌ Произошла внутренняя ошибка.\n⏳ Ожидайте, скоро всё заработает."
)
