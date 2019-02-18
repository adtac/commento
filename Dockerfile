# backend build (api server)
FROM golang:1.10.2-alpine AS api-build

COPY ./api /go/src/commento/api
WORKDIR /go/src/commento/api

RUN apk update && apk add bash make git curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN make prod -j$(($(nproc) + 1))


# frontend build (html, js, css, images)
FROM node:10.3.0-alpine AS frontend-build

COPY ./frontend /commento/frontend/
WORKDIR /commento/frontend/

RUN apk update && apk add bash make
RUN npm install -g html-minifier@3.5.7 uglify-js@3.4.1 sass@1.5.1

RUN make prod -j$(($(nproc) + 1))


# templates build
FROM alpine:3.7 AS templates-build

COPY ./templates /commento/templates
WORKDIR /commento/templates

RUN apk update && apk add bash make

RUN make prod -j$(($(nproc) + 1))


# db build
FROM alpine:3.7 AS db-build

COPY ./db /commento/db
WORKDIR /commento/db

RUN apk update && apk add bash make

RUN make prod -j$(($(nproc) + 1))


# final image
FROM alpine:3.7

COPY --from=api-build /go/src/commento/api/build/prod/commento /commento/commento
COPY --from=frontend-build /commento/frontend/build/prod/*.html /commento/
COPY --from=frontend-build /commento/frontend/build/prod/css/*.css /commento/css/
COPY --from=frontend-build /commento/frontend/build/prod/js/*.js /commento/js/
COPY --from=frontend-build /commento/frontend/build/prod/images/* /commento/images/
COPY --from=frontend-build /commento/frontend/build/prod/fonts/* /commento/fonts/
COPY --from=templates-build /commento/templates/build/prod/templates/ /commento/templates/
COPY --from=db-build /commento/db/build/prod/db/ /commento/db/

RUN apk update && apk add ca-certificates --no-cache

EXPOSE 8080

WORKDIR /commento/

ENV COMMENTO_BIND_ADDRESS="0.0.0.0"
ENTRYPOINT ["/commento/commento"]
