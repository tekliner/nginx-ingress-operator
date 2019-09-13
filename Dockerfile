FROM golang:1.12-alpine

RUN apk update
RUN apk add git mercurial

WORKDIR /app/nginx-ingress-operator
COPY . .

WORKDIR /app/nginx-ingress-operator/cmd/manager

RUN CGO_ENABLED=0 go build
RUN go install

FROM alpine:3.10

COPY --from=0 /go/bin/manager /usr/local/bin/nginx-ingress-operator

ENTRYPOINT ["/usr/local/bin/entrypoint"]
