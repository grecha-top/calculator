# Сервис подсчёта арифметических выражений
___
* Данный интерфейс представляет собой API-интерфейс, реализованный на Go, который позволяет отправлять POST-запросы


> ### Endpoints
> | Эндпоинт | Допустимые методы | Описание |
> | --- | --- | --- |
> | /calculator | *POST* | Получает POST-запрос c телом запроса в формате [JSON](https://ru.wikipedia.org/wiki/JSON), содержащим выражение. Отвечает в формате [JSON](https://ru.wikipedia.org/wiki/JSON) |

---
## Установка:

Введите в терминале команду:
`git clone https://github.com/grecha-top/calculator `

---

## Как использовать:

1. Перейдите в директорию проекта - calculator
2. Введите в командную строку `go run ./cmd/main.go`
3. Отправьте POST-запрос (например, через [curl](https://curl.se/) или [postman](https://www.postman.com/)) на URL: `localhost:8080/calculator`.
4. Получите ответ.

> При возникновении ошибки убедитесь, что установлена версия Go `1.23.3`.
> Последнюю версию можно установить [здесь](https://go.dev/dl/).

## Примеры использования с cURL:

| cURL                                    | Ответ                                     | *HTTP* код
|------------------------------------------------|-------------------------------------------| ----------------------------- |
| ```curl -XPOST -d '{ "expression" : "2 + 2*2"}' 'http://localhost:8080/calculator'```  | ```{"result":6} ``` | 200 |
| ```curl -XPOST -d '{ "expression" : "2 -"}' 'http://localhost:8080/calculator'``` | ```{"error":"invalid expression"}```|400|
| ```curl -XPOST -d '{ "expression" : "2/0"}' 'http://localhost:8080/calculator'``` | ```{"error":"division by zero"}```|400|


