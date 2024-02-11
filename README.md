# go_rest_app

🚀 Simple Go REST App with Telegram Client

This repository contains a straightforward Go application that serves as a RESTful API with integrated Telegram client functionality. The app allows you to interact with Telegram services through a RESTful interface, making it easy to integrate Telegram features into your projects. Whether you're a beginner or an experienced developer, this app provides a clean and minimalistic template to get you started quickly.

## Add credentials to setting/local.yaml file
```sh
vi settings/local.yaml 

-----------
example . fiele

TOKEN=telegram-token
env: local
secret: secret-auth-token
port: 8080
telegram:
    token: telegram-token
    chat_id: chat-id

```
## build
```sh
go build main.go
```
## start
```sh
./main
```