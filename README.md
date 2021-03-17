## Description
dịch vụ này sẽ xử lý nghiệp vụ liên quan tới user

## REQUIREMENT
```
- Mongodb
```

## INSTALL
```bash
#install golang
sudo snap install go --classic
```

## RUN
```bash
go run .
```

## DOC
```bash
// generate doc
//go get -u github.com/swaggo/swag/cmd/swag
//swag init
// access doc
http://localhost/docs/index.html
```

## CONFIG
```.env
NO_SSL_PORT=
ENV=development|production

MONGODB_HOST=host
MONGODB_PORT=27017
MONGODB_USERNAME=
MONGODB_PASSWORD=
```

## Docker
```bash
docker build -t ocr.service.authorization .
```

## TEST (not yet)
unit test
```bash
go test $(go list ./... | grep -v /vendor/ | grep -v /test)
```
e2e test
```bash
go test ./test
```