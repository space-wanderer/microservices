![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/space-wanderer/4b278207b5900240ff8ced6cdb26cf64/raw/coverage.json)

# microservices

Этот репозиторий содержит проект из курса [Микросервисы, как в BigTech 2.0].

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```
Для того чтобы запустить сервисы необходимо выполнить:

```bash
go run inventory/cmd/main.go
go run payment/cmd/main.go
go run order/cmd/main.go
```
