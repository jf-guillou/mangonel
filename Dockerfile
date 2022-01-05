FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY assets/ ./assets/

RUN go build -o /mangonel

VOLUME /app/storage

EXPOSE 8066

CMD [ "/mangonel" ]
