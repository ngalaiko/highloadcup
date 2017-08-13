FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/ngalayko/highloadcup
RUN go get -u github.com/golang/dep/cmd/dep
ADD . ./
RUN make build-alpine

FROM alpine:latest

COPY --from=builder /go/src/github.com/ngalayko/highloadcup/containers/prod/config.yaml .
COPY --from=builder /go/src/github.com/ngalayko/highloadcup/bin/highloadcup .
RUN date > /created_date
EXPOSE 80

CMD ["./highloadcup", "--config=config.yaml"]
