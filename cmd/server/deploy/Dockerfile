FROM golang:latest as builder

WORKDIR /server

COPY go.mod go.sum ./

RUN /usr/local/go/bin/go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux make build-server

FROM alpine:latest

RUN apk update && apk add --no-cache bash
RUN apk add --no-cache tzdata

WORKDIR /server/

COPY --from=builder /server/cmd/server/bin/server .

RUN  chmod +x server
RUN ls

CMD ["./server"]
