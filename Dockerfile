FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /maxscale_exporter

ENTRYPOINT ["/maxscale_exporter -path=/etc/node_exporter/maxscale.json"]
