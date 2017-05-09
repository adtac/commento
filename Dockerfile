# golang1.8.1 base image
FROM golang:1.8.1

# copy workspace direcotie
COPY . /go/src/commento
WORKDIR /go/src/commento

# update package list
RUN apt-get update

# install Node.js and npm
RUN apt-get install -y nodejs npm

# setup nodejs
RUN update-alternatives --install /usr/bin/node node /usr/bin/nodejs 10

# install frontend build dependencies
RUN npm install

# build frontend
RUN npm build

# get dependencies
RUN go get -v -d

# build and install the binary to the $GOPATH/bin
RUN go install

# copy static assets to the binary path
RUN cp -r ./assets /go/bin

# expose port 8080
EXPOSE 8080

# start the server
CMD ["/go/bin/commento"]
