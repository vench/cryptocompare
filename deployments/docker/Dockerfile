ARG GO_VERSION=1.17.2
ARG ALPINE_VERSION=3.14

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
RUN apk add --no-cache \
    curl \
    make \
    g++ \
    gcc
WORKDIR /opt/cryptocompare
COPY . .
ARG ACTION
RUN make $ACTION

FROM alpine:${ALPINE_VERSION}
WORKDIR /opt/cryptocompare
COPY --from=build /opt/cryptocompare/bin/cryptocompare .
CMD ["./cryptocompare"]
