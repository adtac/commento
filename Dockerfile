FROM golang:1.8.1

COPY . /go/src/commento
WORKDIR /go/src/commmento

RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/op/go-logging
RUN go get github.com/adtac/commento
RUN go install github.com/adtac/commento

EXPOSE 8080
CMD ["/go/bin/commento"]
