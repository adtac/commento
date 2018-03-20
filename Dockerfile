# Build backend
FROM golang:1.8.1 as backend-build

COPY . /go/src/commento
WORKDIR /go/src/commento

RUN go get -v .
RUN go build -ldflags '-linkmode external -extldflags -static -w'


# Build frontend
FROM node:8.8-alpine as frontend-build

COPY ./package.json /commento/package.json
COPY ./Gulpfile.js /commento/Gulpfile.js
WORKDIR /commento/

RUN npm install

COPY ./assets/ /commento/assets/

RUN npm run-script build


# Build final image
FROM alpine:3.6

COPY --from=backend-build /go/src/commento/commento /commento/
COPY --from=frontend-build /commento/assets/ /commento/assets/

RUN mkdir /commento-data/
ENV COMMENTO_DATABASE_FILE /commento-data/commento.sqlite3

WORKDIR /commento
ENTRYPOINT /commento/commento
