FROM golang:1.16-alpine3.14 AS builder

WORKDIR /build/
COPY . .
RUN apk add build-base
RUN GOOS=linux
RUN make build

FROM alpine:3.14

RUN adduser --home /home/dorsradek --disabled-password dorsradek
WORKDIR /home/dorsradek
COPY --from=builder /build/app .
USER dorsradek
ENTRYPOINT ["/home/dorsradek/app"]
