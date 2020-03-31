# backend build (api server)
FROM golang:1.14-alpine AS api-build
RUN apk add --no-cache --update bash dep make git curl g++

COPY ./api /go/src/commento/api/
WORKDIR /go/src/commento/api
RUN make prod -j$(($(nproc) + 1))


# frontend build (html, js, css, images)
FROM node:10-alpine AS frontend-build
RUN apk add --no-cache --update bash make python2 g++

COPY ./frontend /commento/frontend
WORKDIR /commento/frontend/
RUN make prod -j$(($(nproc) + 1))


# templates and db build
FROM alpine:3.9 AS templates-db-build
RUN apk add --no-cache --update bash make

COPY ./templates /commento/templates
WORKDIR /commento/templates
RUN make prod -j$(($(nproc) + 1))

COPY ./db /commento/db
WORKDIR /commento/db
RUN make prod -j$(($(nproc) + 1))


# final image
FROM alpine:3.7
RUN apk add --no-cache --update ca-certificates

COPY --from=api-build /go/src/commento/api/build/prod/commento /commento/commento
COPY --from=frontend-build /commento/frontend/build/prod/js /commento/js
COPY --from=frontend-build /commento/frontend/build/prod/css /commento/css
COPY --from=frontend-build /commento/frontend/build/prod/images /commento/images
COPY --from=frontend-build /commento/frontend/build/prod/fonts /commento/fonts
COPY --from=frontend-build /commento/frontend/build/prod/*.html /commento/
COPY --from=templates-db-build /commento/templates/build/prod/templates /commento/templates/
COPY --from=templates-db-build /commento/db/build/prod/db /commento/db/

EXPOSE 8080
WORKDIR /commento/
ENV COMMENTO_BIND_ADDRESS="0.0.0.0"
ENTRYPOINT ["/commento/commento"]
