FROM golang:1.20.1-alpine3.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /go/bin/app

FROM alpine:3.17

WORKDIR /home/app

COPY --from=build /go/bin/app /

EXPOSE 3000

ENTRYPOINT ["/app"]