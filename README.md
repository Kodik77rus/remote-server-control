
- [Реализация](#реализация)
- [Запуск](#запуск)
- [Тесты](#тесты)
- [Api](#api)
  - [POST `https://127.0.0.1:8443/api/v1"/remote-execution`](#post-https1270018443apiv1remote-execution)
      - [Body](#body)
      - [Ошибки](#ошибки)
      - [Пример](#пример)

# Реализация

Сделал адаптивный сервер под linux и под windows, путём добавления флага исполнения к бинарнику с параметром `-path`, параметр принимает путь до папки с sll сертификатом, и адаптирует путь под нужную операционную систему, в моей реализации имя сертификата и ключа зашиты `certificate.pem` и, соответсвенно, `"key.pem"`. Сертификат сгенерировал с помощью `opensll` командой ```openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out certificate.pem```, находиться в папке `configs`.

>Пытался добиться идеальной работы на windows, но так, как у меня стоит пакет с русским языком, у меня полетела кодировка, пытался парсить в другие кодировки, у меня не получилось, но на английском пакете, по идее, всё нормально работает.

При написании данного сервера старался  придерживаться следующих ресурсов:
```https://github.com/sau00/uber-go-guide-ru/blob/master/style.md```</br>
```https://github.com/golang-standards/project-layout```

# Запуск

под Linux
```Makefile
 make build
```

под Windows надо указать путь

# Тесты

запуск
```Makefile 
 make test
```
>Ужасные тесты написал, времени не хватило нормально написать, они корректно работают только с линуксом. CI/CD не добавил.

# Api

REST архитектура

## POST `https://127.0.0.1:8443/api/v1"/remote-execution`

Возвращает json в такой структуре

```json
{
    "message": {
        "stderr": "",
        "stdout": ""
    }
}

```

#### Body

| Полe json |   Тип    |         Описание         | Важность параметра |
| --------- | :------: | :----------------------: | :----------------: |
| **cmd**   | `string` | имя исполняемой комманды |        `Да`        |
| **os**    | `string` | тип операциооной системы |        `Да`        |
| **stdin** | `string` | входной поток в комманду |       `Нет`        |

#### Ошибки

| Название                      |  Код  | Описание                                        |
| :---------------------------- | :---: | :---------------------------------------------- |
| interrupted                   |  400  | Долгое исполнение команды, лимит 5 сек          |
| wrong os type                 |  400  | Неправильно переданный тип операционой  системы |
| failed execute command        |  400  | Ошибка запуска программы                        |
| empty command name            |  400  | В теле запроса пустое поле `имя команды`        |
| wrong type provided for field |  400  | Неверный тип в теле запроса                     |
| unknown field                 |  400  | Неизвестное поле в теле запроса                 |

#### Пример

Отправляем

```json
{
    "cmd": "tr a-z A-Z", 
    "os": "linux",
    "stdin": "test"
}
```

Получаем

```json
{
    "message": {
        "stderr": "",
        "stdout": "TEST"
    }
}
```
