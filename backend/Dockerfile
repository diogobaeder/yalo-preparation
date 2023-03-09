FROM golang:1.20

RUN mkdir /app
WORKDIR /app
COPY . /app/
RUN go mod download && go mod verify

RUN go build -v