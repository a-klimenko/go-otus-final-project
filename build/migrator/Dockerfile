FROM golang:1.17

WORKDIR "/migrations"

COPY ./migrations .

RUN go get -u github.com/pressly/goose/cmd/goose