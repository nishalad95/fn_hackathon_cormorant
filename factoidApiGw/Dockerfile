FROM golang:1.11-alpine

WORKDIR /go/src/factoidApiGw
COPY main.go .

RUN go install -v ./...

CMD ["factoidApiGw"]
