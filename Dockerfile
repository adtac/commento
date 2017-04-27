FROM golang:1.8.1

RUN mkdir /go/src/commento
ADD . /go/src/commento/
WORKDIR /go/src/commento

RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/op/go-logging
RUN go build -o commento .

EXPOSE 8080
CMD ["/go/src/commento/commento"]
