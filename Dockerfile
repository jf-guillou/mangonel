FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /mangonel

FROM alpine:3.15

WORKDIR /app

COPY --from=build /mangonel /app/mangonel
COPY assets/ /app/assets/

VOLUME /app/storage

EXPOSE 8066

CMD [ "/app/mangonel" ]
