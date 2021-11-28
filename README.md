# Wikidump Import to Mongodb With Golang
## Prerequisite 

Sebelum memulai, silahkan mengganti username dan password pada file main.go line berikut

```
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://username:password@localhost:27017"))
```
kemudian ganti nama database (disini 'idwiki') dan collection (disini 'content') pada file main.go line berikut jika diperlukan

```
collections := client.Database("idwiki").Collection("content")
```
kemudian ganti nama file (disini 'dump.json') pada file main.go line berikut jika diperlukan

```
f, err := os.Open("dump.json")
```
## Usage

Cukup run command berikut di terminal. Selanjutnya akan muncul id wikicontent yang telah dibuat pada terminal

```
go run .\main.go
```
