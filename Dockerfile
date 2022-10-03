FROM golang:1.19.1-alpine3.16 as builder
ENV GOPROXY https://goproxy.cn
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -v -o webhook .

FROM alpine:3.16
RUN apk add bash
COPY --from=builder /build/webhook /usr/local/bin
COPY ./entrypoint.sh /usr/local/bin
ENTRYPOINT ["entrypoint.sh"]
