[[EN]](../README.md) **[RU]**

# Telegram-processor

Инструмент для исследования семантического поиска по сообщениям в Telegram

---

## Как пользоваться?

### Скачать

https
```bash
git clone https://github.com/jesus-development/telegram-processor.git && cd telegram-processor
```

или ssh
```bash
git clone git@github.com:jesus-development/telegram-processor.git && cd telegram-processor
```

### Подготовка

#### Docker (рекомендуется)

```bash
cp -n .env.example .env
```

#### Локально

```bash
make build-demo
```

### Добавить свой OpenAI Api Token (optional)

Добавьте `OPENAI_API_KEY=ваш_токен` в `.env` если хотите.  
Иначе векторы и результаты поиска будут случайными.

### Проверить порты (optional)

Тестовые контейнеры используют внешние порты:

- `5432` - postgres
- `50051` и `50052` - grpc и http api server

Если они уже заняты:

- Изменить их в `docker-compose.yml`
- Если хочется запустить демо локально, также изменить `db.port` в `configs/local.yaml`


### Запустить контейнеры с БД, сервером и демонстрацией

```bash
docker-compose up -d
```

### Запустить демонстрацию

#### Docker (рекомендуется)

  При первом запуске используйте флаг `--import-db` чтобы импортировать тестовые 30 сообщений из `resources/demo/chat-export-30news.json`.
```bash
docker exec -it telegram-processor_demo ./telegram-processor demo --import-db
```

#### Локально
Требуется go 1.22  

При первом запуске используйте флаг `--import-db` чтобы импортировать тестовые 30 сообщений из `resources/demo/chat-export-30news.json`.

```bash
./telegram-processor demo -c configs/local.yaml --import-db
```

### Отправить запрос

```bash
curl --location '127.0.0.1:50052/api/search?query=asdf'
```

### Остановить контейнеры

```bash
docker-compose down
```

---

## Documentation

### Project structure

- `cmd` - cli commands
-
    - `root` - base for all commands
-
    - `demo` - demonstration
-
    - `embeddings` - embedding tools
-
    - `import-json` - import messages from json
-
    - `server` - gRPC and http server
- `configs` - config files
- `google`, `grpc-gateway` - static proto libraries
- `internal` - main logic
- - `api` - api server, handlers
- - `config` - config structs
- - `db` - database connections
- - `logger` - log customization
- - `models` - entities
- - `repository` - repositories
- - `scenarios` - scenarios for demo and tests
- - `services` - modules with business logic
- `pkg` - utils, proto/json models, etc.
- `resources` - csv/json-files, images, static docs, etc.
- `scripts`
- - `bash/wait-for-it.sh` - script for waiting db connection in docker
- - `db/init.sql` - init sql script
### Stack

- Embeddings: [OpenAI](https://platform.openai.com/docs/guides/embeddings)
- DB: [Postgres/pgVector](https://github.com/pgvector/pgvector)

---

## FAQ
- **Почему результаты поиска такие плохие?**  
Смысл всего текста сжимается в один вектор. 
Из-за чего отдельные части текста мешают найти друг друга. 
- **Что делать?**  
Ждать следующего релиза с предобработкой текста и гибридным поиском.