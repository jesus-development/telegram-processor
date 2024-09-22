**[EN]** [[RU]](doc/README_ru.md)

# Telegram-processor

Research tool for semantic search using Telegram messages

---

## Getting started

### Clone

https
```bash
git clone https://github.com/jesus-development/telegram-processor.git && cd telegram-processor
```

or ssh
```bash
git clone git@github.com:jesus-development/telegram-processor.git && cd telegram-processor
```

### Init

#### Docker (recommended)

```bash
cp -n .env.example .env
```

#### Locally

You need go version 1.22
```bash
make build-demo
```

### Set your OpenAI Api Token (optional)

Add `OPENAI_API_KEY=your_token` to `.env` if you have it.  
Otherwise embeddings and search result will be random.

### Check ports (optional)

Test containers use external ports:

- `5432` - postgres
- `50051` and `50052` - grpc and http api server

If they already in use:

- Change them in `docker-compose.yml`
- If you want to run demo locally, also change `db.port` in `configs/local.yaml`


### Run containers with postgres and api server

```bash
docker-compose up -d
```

### Run demo

#### Docker (recommended)

  For the first run, use `--import-db` for import test messages from `resources/demo/chat-export-30news.json`.
```bash
docker exec -it telegram-processor_demo ./telegram-processor demo --import-db
```

#### Locally

For the first run, use `--import-db` for import test messages from `resources/demo/chat-export-30news.json`.

```bash
./telegram-processor demo -c configs/local.yaml --import-db
```

### Send API request directly

```bash
curl --location '127.0.0.1:50052/api/search?query=asdf'
```

### Stop containers

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
- **Why search results so bad?**  
The meaning of large text is compressed into a single vector. 
Therefore, the individual components of the text reduce the accuracy of retrieving each other. 
- **And what to do?**  
Wait for the next release with text preprocessing and hybrid search.