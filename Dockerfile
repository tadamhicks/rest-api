FROM golang:1.20.10

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["./rest-api"]