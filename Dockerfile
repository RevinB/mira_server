FROM golang:1.18.0-alpine AS build
RUN apk add git

ARG VERSION=dev

WORKDIR /tmp/app
COPY . .

RUN go mod download && \
    go mod verify && \
    go build -ldflags="-X 'main.version=${VERSION}'" -o entry .

FROM alpine:latest

WORKDIR /app
COPY --from=build /tmp/app/entry /app/entry

EXPOSE 3000

ENTRYPOINT ["./entry"]