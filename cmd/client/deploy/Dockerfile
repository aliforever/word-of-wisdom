FROM golang:latest as builder

WORKDIR /client

COPY go.mod go.sum ./

RUN /usr/local/go/bin/go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux make build-client

FROM alpine:latest

RUN apk update && apk add --no-cache bash
RUN apk add --no-cache tzdata

WORKDIR /client/

COPY --from=builder /client/cmd/client/bin/client .

RUN  chmod +x client
RUN ls

CMD ["./client"]
