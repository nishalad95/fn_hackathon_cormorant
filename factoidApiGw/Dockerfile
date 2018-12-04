FROM golang:1.11-alpine

WORKDIR /go/src/factoidApiGw
COPY . .

RUN go install -v ./...

CMD ["factoidApiGw"]
