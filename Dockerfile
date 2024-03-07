ARG GOLANG_VERSION=1.22-alpine
FROM golang:${GOLANG_VERSION} AS builder
MAINTAINER Alexandre Ferland <me@alexferl.com>

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" ./cmd/server

FROM scratch
COPY --from=builder /build/server /server
COPY --from=builder /build/configs /configs

ENTRYPOINT ["/server"]

EXPOSE 1323
