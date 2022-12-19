# Gophkeeper
Выпускной проект для курса "Продвинутый разработчик Go" от Яндекс практикум

## Общее описание
Система представляет собой менеджер паролей: она обеспечивает безопасное хранение и передачу следующих данных пользователя:
* пара логин/пароль
* данные кредитных карт
* Произвольные текстовые строки
* произвольные файлы (до 1 Мбайта)

### Клиентский интерфейс
![demo](docs/gk-demo.gif)

## Состав системы
Система представляет собой клиент серверное приложение состоящее из
* клиетской программы (client)
* серверной программы (server)
* базы данных (db)

### Протокол взаимодействия
Клиент и сервер обмениваются данными с помощью gRPC с использованием зашифрованного TLS туннеля.


### Client
Представляет интерактивную программу с интерфейсом коммандной строки.
Клиент реализует следующий функционал:
* регистрация нового пользователя
* аутентификация существующего пользователя
* щифрование, расшифровка приватных данных пользователя
* создание, удаление обновление приватных данных пользователя

### Server
Сервер обслуживает запросы от клиентской программы и реализует следующую логику:
* обслуживание запроса на регистрацию
* аутентификация пользователя
* обслуживание запросов на создание/обновление/удаление приватных данных пользователя
* работа с базой данных

### DB
Осуществляет долгострочное хранение данных системы: данные о зарегистрированныхп ользователей, а также их приватная информация.


## Модель защиты

### Регистрация
Пароли пользователей храняться в БД в зашифрованном виде. Шифрование происходит с помощью bcrypt.
|![register.svg](docs/register.svg)|
|:--:|
| *регистрация пользователя* |


### Аутентификация и авторизация
Для аутентификации запросов пользователя, используются JWT токены. Токен генерится при аутентификации пользователя и отправляются со всеми командами (кроме register/login). Токен (как и вся информация между сторонами) шифруется с помощью TLS
|![login.svg](docs/gophkeeper-login.drawio.svg)|
|:--:|
| *аутентификация пользователя / авторизация запроса* |

### Шифрование данных
Пользовательские данные шифруются с помощью мастер ключа, который хранится только на клиенте. Алгоритм шифрования AES. Вместе с зашифрованными данными на сервер отправляется хеш сумма ключа. Это позволит в будущем определить каким ключем были зашифрованы данные.
> **Важно:**  Мастер ключ хранится в незашифрованном виде. Защита ключа возлагается на пользователя.

|![gophgophkeeper-encrypt.drawio.svg](docs/gophkeeper-encrypt.drawio.svg)|
|:--:|
| *шифрование/расшифровка данных* |


## Конфигурация системы

### Клиент
| Параметр                 | Флаг             | ENV                               | Обязательное | Значение по умолчанию               |
|--------------------------|------------------|-----------------------------------|--------------|-------------------------------------|
| Адрес и порт сервера     | -a localhost:443 | GK_SERVER_ADDRESS="localhost:443" | Нет          | localhost:443                       |
| Использовать  TLS (bool) | -t               | N/A                               | Нет          | false (шифрование не используется)  |
> **Важно:**  если задан и флаг и переменная окружения, приоритет за переменной окружения.


### Пример запуска клиента
Ниже представлен пример запуска клиента который будет подключаться к серверу exemple.com по порту 443 с использованием TLS
```bash
cmd/client/client -t -a example.com:443
```

### Сервер
| Параметр                              | Флаг                                                         | ENV                                                                   | Обязательное | Значение по умолчанию                         |
|---------------------------------------|--------------------------------------------------------------|-----------------------------------------------------------------------|--------------|-----------------------------------------------|
| Адрес и порт сервера                  | -a localhost:443                                             | GK_SERVER_ADDRESS="localhost:443"                                     | Нет          | localhost:443                                 |
| DSN базы данных                       | -d postgres://uname:pass@localhost:5432/db | GK_DB_DSN="postgres://uname:pass@localhost:5432/db" | Да           | N/A                                           |
| Ключ для подписи JWT токена           | -k secret                                                    | GK_KEY="secret"                                                       | Да           | N/A                                           |
| Путь к сертификату сервера (.pem)     | -c ./certs/cert.pem                                          | GK_CERT="./certs/cert.pem"                                            | Нет          | Если не задано, сервер будет работать без TLS |
| Путь к закрытому ключу сервера (.pem) | -p ./certs/key.pem                                           | GK_PRIVATE_KEY="./certs/key.pem"                                      | Нет          | Если не задано, сервер будет работать без TLS |
> **Важно:**  если задан и флаг и переменная окружения, приоритет за переменной окружения.

### Пример запуска сервера
Ниже представлен пример запуска сервера на сокете localhost:10443 с использованием TLS и подключением к локальной базе данных gophkeeper
```
cmd/server/server -d postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper -a localhost:10443 -k secret -c ./certs/server-cert.pem -p ./certs/server-key.pem
```

## Локальный деплой
Ниже представлена инструкция о том как развернуть проект в локальном окружении (для разработки или тестов)


### Требования
* docker + docker-compose
* утилита make
* openssl
* git

### Инструкция
1. Клонировать проект
    ```bash
    [kzhukov@fedora gophkeeper]$ git clone git@github.com:zklevsha/gophkeeper.git && cd ./gophkeeper
    [.. output ommited ...]

    ```
2. Cобрать клиент и сервер
    ```bash
    [kzhukov@fedora gophkeeper]$ make client && make server
    cd ./cmd/client/ && go build .
    cd ./cmd/server/ && go build .
    [kzhukov@fedora gophkeeper]$

    ```
3. Сгенерить сертификат и ключ для сервера
    ```bash
    [kzhukov@fedora gophkeeper]$ make certs
    mkdir certs && \
    echo "Generating CA cert and key" && \
    openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout certs/ca-key.pem -out certs/ca-cert.pem -subj "/C=RU/L=Moscow/O=Practicum/OU=Practicum/CN=CA server" && \
    echo "Generate server key and sign request" && \
    openssl req -newkey rsa:4096 -nodes -keyout certs/server-key.pem -out certs/server-req.pem -subj "/C=RU/L=Moscow/O=Practicum/OU=Practicum/CN=CA server" && \
    echo "Generating servers cert" && \
    openssl x509 -req -in certs/server-req.pem -days 120 -CA certs/ca-cert.pem -CAkey certs/ca-key.pem -CAcreateserial -out certs/server-cert.pem

    [...output ommited...]
    ```

3. Запустить контейнеры баз данных
    ```bash
    [kzhukov@fedora gophkeeper]$ sudo docker-compose up -d
    Recreating db      ... done
    Recreating db_test ... done
    [kzhukov@fedora gophkeeper]$
    ```
    ( db используется непосредственно сервером, db_test используется для юнит тестов)

4. Запустить миграции баз данных
    ```bash
    kzhukov@fedora gophkeeper]$ source .env && make migrate_up
    migrate -database postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable -path db/migrations up
    1/u init (39.561859ms)
    2/u private_data_add_indexes (62.789628ms)
    3/u private_data_fix_unique (83.884704ms)
    4/u add_card_pdata_type (99.50373ms)
    5/u private_types_add_pstring_type (114.56904ms)
    6/u private_types_add_pfile_type (129.026149ms)
    ```

5. Запустить сервер
    ```bash
    [kzhukov@fedora gophkeeper]$ ./start_server.sh
    2022/12/19 15:12:04 starting server on localhost:10443
    ```

6. Запустить клиента
    ```bash
    [kzhukov@fedora gophkeeper]$ ./start_client.sh
    Welcome to gophkeeper. Let`s set you up
    ✔ register
    Registering:
    email: t@t.ru
    password: ****
    password(confirm): ****
    Register succsessful.
    ```


## Ссылки
 - [Teхническое задание](docs/terms-of-reference.md)
