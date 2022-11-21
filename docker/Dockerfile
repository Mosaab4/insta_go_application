FROM golang:1.18-alpine

RUN mkdir /app

WORKDIR /app

COPY ../go.mod ./

RUN go mod download

COPY ../ /app

RUN env GOOS=linux CGO_ENABLED=0 go build -o goApplication ./

CMD ["/app/goApplication"]