# `goph-keeper-client`

## CLI утилита клиент сервиса goph-keeper 

# Первый запуск
* запустить `./goph-keeper-client -gen-config`
* запустить `./goph-keeper-client -register -username "your user name" -password "your password"`

## строка запуска `./goph-keeper-client [command] [options]`
### `ВНИМАНИЕ!`Можно передавать только 1 команду

### commands

#### первый запуск
* `-gen-config`        сгенерировать конфиг
* `-register`          регистрация в сервисе
#### профиль пользователя
* `-profile`           информация о профиле
* `-gen-keys`          сгенерировать новую пару RSA ключей
* `-change-password`   сменить пароль
#### данные
* `-get`               получить данные
* `-put`               сохранить данные
* `-remove`            удалить данные
* `-list`              список данных 

### common options
* `-config`            файл конфига
* `-not-store-config`  не сохранять конфиг

### `gen-config`

Сгенерировать пустой файл конфига в формате JSON

#### options
* `-path`          (опционально) путь к файлу конфига

### `register`

Регистрация пользователя в сервисе goph-keeper

#### options
* `-address`       хост:порт (http://example.org:8080)
* `-username`      имя пользователя
* `-password`      пароль
* `-person`        (опционально) имя и фамилия
* `-e-mail`        (опционально) e-mail

### `profile`

Вывести информацию о профиле

#### options
* `-address`       (опционально) хост:порт (http://example.org:8080)
* `-username`      (опционально) имя пользователя
* `-password`      (опционально) пароль

### `gen-keys`

Сгенерировать новую пару RSA ключей, результат будет сохранён в принятый конфиг файл

#### options
* `-address`       (опционально) хост:порт (http://example.org:8080)
* `-username`      (опционально) имя пользователя
* `-password`      пароль

### `change-password`

Произвести смену пароля

#### options
* `-address`       (опционально) хост:порт (http://example.org:8080)
* `-username`      (опционально) имя пользователя
* `-password`      (опционально) пароль
* `-old-password`  старый пароль, обязан совпадать с текущим паролем
* `-new-password`  новый пароль, не пустой, не совпадает со старым

### `get`

Получить сохранённые ранее данные

#### options
* `-address`       (опционально) хост:порт (http://example.org:8080)
* `-username`      (опционально) имя пользователя
* `-password`      (опционально) пароль
* `-data-kind`     тип сохранённых данных (credential,plaintext,binary,bankcard)
* `-name`          наименование данных
* `-path`          (опционально) полный путь к файлу (по умолчанию в каталог приложения)

### `put`

Сохранить данные

#### option
* `-address`       (опционально) хост:порт (http://example.org:8080)
* `-username`      (опционально) имя пользователя
* `-password`      (опционально) пароль
* `-data-kind`     тип сохранённых данных (credential,plaintext,binary,bankcard)
* `-name`          наименование данных
* `-data`          текстовые данные
* `-path`          полный путь к файлу, это опция необходима при типе данных binary

### `remove`

Удалить ранее сохранённые данные

#### options
* `-address`       (опционально) хост:порт (http://example.org:8080)
* `-username`      (опционально) имя пользователя
* `-password`      (опционально) пароль
* `-data-kind`     тип сохранённых данных (credential,plaintext,binary,bankcard)
* `-name`          наименование данных

### list

Получить список сохранённых данных

#### options
* `-address`       (опционально) хост:порт (http://example.org:8080)
* `-username`      (опционально) имя пользователя
* `-password`      (опционально) пароль

### Что означает (опционально) в описании опций? Ваши credentials записываются в конфиг. Если эти опции не указаны, значит берутся из конфиг файла