FROM golang:1.13-alpine as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev bash linux-headers eudev-dev

WORKDIR /hackcli

# Add source files
COPY . .

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

RUN make install

# ----------------------------

FROM alpine

COPY --from=builder /go/bin/hackcli /usr/local/bin/hackcli
