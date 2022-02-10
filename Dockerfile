FROM golang:1.17-alpine3.15

RUN apk add --no-cache git

WORKDIR /app/lostarkstatus

COPY go.mod .

RUN go mod download

COPY . .

ENTRYPOINT ["go", "run", "main.go"]
