FROM golang:1.12-alpine

WORKDIR /go/src/github.com/tekliner/nginx-ingress-operator
COPY . .

WORKDIR /go/src/github.com/tekliner/nginx-ingress-operator/cmd/manager

RUN CGO_ENABLED=0 go build
RUN go install

FROM alpine:3.9

COPY --from=0 /go/bin/manager /usr/local/bin/nginx-ingress-operator

ENTRYPOINT ["/usr/local/bin/entrypoint"]
