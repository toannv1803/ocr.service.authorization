FROM golang:alpine3.12 AS builder
WORKDIR /go/src/app
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/ocr.service.backend

FROM alpine:3.12
WORKDIR /go/src
COPY --from=builder /go/bin/ocr.service.backend /go/src/ocr.service.backend
COPY --from=builder /go/src/app/docs /go/src/docs
CMD ["/go/src/ocr.service.backend"]