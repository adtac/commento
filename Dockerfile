# golang1.8.1 base image
FROM golang:1.8.1

# copy source
COPY . /go/src/commento
WORKDIR /go/src/commento

# build backend
RUN go get -v .
RUN go install .

# build frontend
RUN git clone https://github.com/creationix/nvm.git /tmp/nvm
RUN /bin/bash -c "source /tmp/nvm/nvm.sh && nvm install node"
RUN /bin/bash -c "source /tmp/nvm/nvm.sh && npm install"
RUN /bin/bash -c "source /tmp/nvm/nvm.sh && npm run-script build"
RUN cp /go/src/commento/assets /go/bin/assets -vr

# set entrypoint
ENTRYPOINT /go/bin/commento
