FROM golang:1.16-alpine3.14 AS builder

WORKDIR /build/
COPY . .
RUN apk add build-base
RUN GOOS=linux
RUN make build

FROM alpine:3.14

RUN adduser --home /home/radek --disabled-password radek
WORKDIR /home/radek
COPY --from=builder /build/app .
USER radek
ENTRYPOINT ["/home/radek/app"]
