# -subscription_service

**Тестовое задание Junior Golang Developer — Effective Mobile**

REST API серсив для управления подписками пользователей и подсчёта суммарной стоимости за период.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)

## Выполненные требования

- HTTP-эндпоинты для CRUDL-операций над подписками
- Подсчёт суммарной стоимости за период с фильтрацией по пользователю и сервису
- PostgreSQL + миграции (goose)
- Логирование запросов и ошибок
- Конфигурация через `.env`
- Полная интерактивная Swagger-документация
- Запуск через `docker-compose up`
- Дополнительно: `/health-check` для мониторинга состояния сервиса и БД

## Технологии

- **Go** + **chi** (роутер)
- **Go "testing"** - для юнит и интеграционных тестов
- **PostgreSQL** - основная БД
- **goose** — миграции БД
- **swaggo/swag** + **http-swagger** — документация
- **Docker** 
- **Makefile** - для сборки проекта

### 📋 **Требования:**
- Git
- Go 1.24.5
- Make (опционально)
- [Goose](https://github.com/pressly/goose)
- [Swagger](https://github.com/swaggo/swag)
- [Docker](https://www.docker.com/)

## 🚀 Быстрый старт

## Быстрый запуск в контейнере (рекомендуемый способ)

```bash
git clone https://github.com/Piccadilly98/subscription_service.git
cd subscription_service
cp .env.example .env
docker-compose up --build
```
[Информация по сервису](#поздравляю-сервер-запущен-и-успешно-работает)

### 🔧 **Установка и запуск:**

#### **Шаг 1: Клонирование проекта**
```bash
git clone https://github.com/Piccadilly98/subscription_service.git
cd subscription_service
```     

#### **Шаг 2: Установка зависимостей:**
```bash
go mod download #для установки всех пакетов требующихся в проекте
go install github.com/swaggo/swag/cmd/swag@latest #для сборки документаций
go install github.com/pressly/goose/v3/cmd/goose@latest #для миграций
```


#### **Шаг 3: Создать конфигурационный файл .env в корне репозитория для запуска сервера как в примере файла [.env.example](https://github.com/Piccadilly98/subscription_service/blob/pre-relese/.env.example) :**
```bash
DB_HOST=localhost           #адрес хоста где расположена БД, дефолтное значение localhost
DB_SSLMODE=disable          #ssl mode, дефолтное значение disable
DB_NAME=subscriptions       #имя базы данных в которой наши таблицы
DB_PASSWORD=postgres        #пароль от юзера
DB_PORT=5432                #порт базы данных, дефолтное значение 5432
DB_USER=postgres            #пользователь базы данных

CACHE=true                  #нужен ли в нашем сервисе кэшm, возможные значения: "y", "yes", "true", "t", "1", DEFAULT false
CACHE_TTL_IN_SECONDS=300    #ttl нашего кэша значения могут быть: integer > 0 or empty, DEFAULT 300
LOGGING_USER_ERROR=false     #хотим ли мы логировать пользовательские ошибки в запросах, возможные значения: "y", "yes", "true", "t", "1", DEFAULT false

SERVER_ADDR=localhost       #адрес на котором запускается наш сервис, DEFAULT "localhost"
SERVER_PORT=8080            #порт на котором мы хотим запускать наш сервис ONLY INTEGER > 0, DEFAULT 8080

# переменные для миграций

GOOSE_DRIVER=postgres                   #какую базу данных мы используем   
GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/subscriptions   #строчка для подключения к нашей базе данных в формате: название_субд://имя_пользователя:пароль@хост:порт/имя базы данных
GOOSE_MIGRATION_DIR=./migrations        #место где хранятся файлы миграций
```

#### **Шаг 4(оптимальный): Запуск тестирования:**
**Если есть make:**
```bash
make tests-all #для запуска тестирования с логами + покрытием
make tests-cover #для запуска тестов только с покрытием
make tests-v #для запуска тестирования только с логами
```
**Если нет make:**
```bash
go test -v -cover ./... #для запуска тестирования с логами + покрытием
go test -cover ./... #для запуска тестов только с покрытием
go test -v ./... #для запуска тестирования только с логами
```

### **Локальная сборка**
**Если есть make:**
```bash
    swag init --dir cmd,internal --generalInfo main/main.go # для сборки документации swagger
    make #для примения миграций, сборки и запуска сервиса
```
[Информация по сервису](#поздравляю-сервер-запущен-и-успешно-работает)




**Если нету make:**
```bash
   swag init --dir cmd,internal --generalInfo main/main.go # для сборки документации swagger
   goose up #для поднятия миграций в нашу базу данных
   go build -o main ./cmd/main #для сборки проекта
   ./main #для запуска бинарника
```

[Информация по сервису](#поздравляю-сервер-запущен-и-успешно-работает)

## **🎉Поздравляю! Сервер запущен и успешно работает**   
- **Слушает порт: 8080 если вы не меняли его в конфиге**
- **Расположен на localhost (или на адресе который вы указали в конфиге)**
- **База данных в случае поднятия в контейнере работает на localhost:5433(что бы не занимать стандартный порт PostgreSQL)**
- **В случае удачного запуска мы должны увидеть подобные логи:**
```bash
    2025/12/31 11:55:30 Database connected successfully
    [SERVER INFO]2025/12/31 11:55:30 INFO: server start in 0.0.0.0:8080
```
- В случае ошибки мы видим либо логи goose, например:
```bash
    2025/12/31 15:56:09 goose run: failed to connect to `user=postgres database=subscriptions` #значит что мы указали ссылку для подключения к бд неверно
```
**или:**
```bash
2025/12/31 15:57:40 goose run: ./migration directory does not exist #значит что мы указали неверный путь к папке с миграциями
 ```

 **Или логи от нашего сервиса:**
 ```bash
 [DB CONNECTION ERROR] 2025/12/31 16:15:46 CRITICAL: error in ping dataBase: dial tcp: lookup 127.128.22.258: no such host #при инвалидном хосте
 ```
- [Документация](#документация)

## **Документация:**
- Документация swagger доступна по адресу:[localhost:8080/swagger/](http://localhost:8080/swagger/)  

- Общая документация по всем хендлерам
<img src="./screenshots/Screenshot 2025-12-31 at 21.40.45.png" width="800" alt="Swagger документация"> 
- Каждый хендлер закументирован следующим образом:
<img src="./screenshots/Screenshot 2025-12-31 at 21.47.50.png" width="800" alt="Swagger хендлер"> 
   - Есть информация о телах которые принимаются и отдаются при каждом статус коде
   <img src="./screenshots/Screenshot 2025-12-31 at 21.51.59.png" width="800" alt="Swagger хендлер подробнее"> 


## Особенности реализации

- Автоматическая валидация и нормализация дат ("MM-YYYY" ↔ "DD-MM-YYYY")
- Динамический статус подписки (active / ended / not started)
- [In-memory кэш](#кэшированиеå) для оптимизации проверки существования подписки по ID
- Консистентные JSON-ошибки через единую структуру
- Multi-stage Dockerfile для минимального финального образа
- Миграции применяются автоматически при старте

## **Технические детали**
### **Слои сервиса:**
**Сервис разделён на 3 основных слоя:**
- БД методы - слой который только делает запросы в базу данных и возвращает целевые данные и ошибку, реализованы 2 вида методов: с использованием транзакций, для избежания гонки данных при изменениях в бд и методы в которые работают в единой транзакции
- Сервис - методы которые обращаются к БД методам и содержат бизнес-логику   
- Хендлеры - обработчики, которые валидируют входные данные, совершают анмаршалинг и совершают запросы к слою сервиса   

### **Кэширование:**
**Для оптимизации запросов к бд в хендлерах где требуется ID подписки в URL был сделан in-memory кэш:**   
- Кэш можно включать с помощью переменной окружения в .env файле:  
```bash
CACHE=true                  #may be "y", "yes", "true", "t", "1", DEFAULT false
```
Дефолтное значение **false**. 

- Настроить ttl в секундах можно с помощью перменной окружения:   
```bash
CACHE_TTL_IN_SECONDS=300    #may be only integer > 0 or empty, DEFAULT 300
```
Дефолтное значение ttl - 300 секунд (5 минут)  


**В коде реализация кэша выглядит так:**     
```go
type Cache struct {
	cache map[string]time.Time   // map хранящяя в себе ID и дату создания записи в кэше
	mu    sync.RWMutex          // mutex для конкурентного доступа к мапе
	ttl   time.Duration         // время жизни записи после истечения, которого запись удаляется из кэша
}
```
**В методах сервиса где требуется id в URL:**  
При значении переменной окружения **CACHE=true:**    
- Сначала проверяется кэш и если id есть в кэше то запрос в базу данных не совершается, для проверки существования этого id
- Если вдруг в кэше данный id есть но у него истёк ttl то кэш, удаляет эту запись и дальше сервис проверяет существование id с помощью запроса в базу данных
- Если записи нету, то сервис проверяет базу данных и в случае существования данного id сохраняет запись в кэше с текущим временем

**При обновлении данных подписки:**
- Сервис обновляет время у записи в кэше


### **БД в контейнере:**
- При запуске с помощью:
  ```bash
  docker-compose up
  ```
  В корне репозитория создаётся директория: **pg_data** которая прокинута в контейнер и при перезапуске контейнера данные сохраняются


## Спасибо за интересное тестовое задание!  
Буду рад любому фидбеку!  

## Контакты:   
- [Резюме на hh](https://samara.hh.ru/resume/875f5df3ff0fbfea5f0039ed1f39506e574431)
- [Резюме на гугл диске](https://drive.google.com/file/d/1I3B2MahiVC0Lcrtt1TexjlfFCZZaV2uU/view?usp=sharing)
- [GitHub](https://github.com/Piccadilly98)