FROM golang:1.22

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY pkg/*.go ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build

EXPOSE 8080

CMD ["./ramen"]
