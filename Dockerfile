FROM golang:1.12.6-alpine as builder
MAINTAINER Alexandre Ferland <aferlandqc@gmail.com>

ENV GO111MODULE=on

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /build/echo-boilerplate /echo-boilerplate

ENTRYPOINT ["/echo-boilerplate"]

EXPOSE 1323
CMD ["--address", "0.0.0.0:1323"]
