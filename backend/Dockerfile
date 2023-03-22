FROM golang:1.20

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download && go mod verify

COPY . /app/
RUN go mod download && go mod verify

RUN go build -o bin/ ./...