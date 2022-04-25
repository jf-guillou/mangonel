FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY assets/ ./assets/

RUN go build -o /mangonel

FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=build /mangonel /mangonel

VOLUME /app/storage

EXPOSE 8066

USER nonroot:nonroot

CMD [ "/mangonel" ]
