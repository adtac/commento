# backend build (api server)
FROM golang:1.15-alpine AS api-build
RUN apk add --no-cache --update bash dep make git curl g++

ARG RELEASE=prod
COPY ./api /go/src/commento/api/
WORKDIR /go/src/commento/api
RUN make ${RELEASE} -j$(($(nproc) + 1))


# frontend build (html, js, css, images)
FROM node:12-alpine AS frontend-build
RUN apk add --no-cache --update bash make python2 g++

ARG RELEASE=prod
COPY ./frontend /commento/frontend
WORKDIR /commento/frontend/
RUN make ${RELEASE} -j$(($(nproc) + 1))


# templates and db build
FROM alpine:3.13 AS templates-db-build
RUN apk add --no-cache --update bash make

ARG RELEASE=prod
COPY ./templates /commento/templates
WORKDIR /commento/templates
RUN make ${RELEASE} -j$(($(nproc) + 1))

COPY ./db /commento/db
WORKDIR /commento/db
RUN make ${RELEASE} -j$(($(nproc) + 1))


# final image
FROM alpine:3.13
RUN apk add --no-cache --update ca-certificates

ARG RELEASE=prod

COPY --from=api-build /go/src/commento/api/build/${RELEASE}/commento /commento/commento
COPY --from=frontend-build /commento/frontend/build/${RELEASE}/js /commento/js
COPY --from=frontend-build /commento/frontend/build/${RELEASE}/css /commento/css
COPY --from=frontend-build /commento/frontend/build/${RELEASE}/images /commento/images
COPY --from=frontend-build /commento/frontend/build/${RELEASE}/fonts /commento/fonts
COPY --from=frontend-build /commento/frontend/build/${RELEASE}/*.html /commento/
COPY --from=templates-db-build /commento/templates/build/${RELEASE}/templates /commento/templates/
COPY --from=templates-db-build /commento/db/build/${RELEASE}/db /commento/db/

EXPOSE 8080
WORKDIR /commento/
ENV COMMENTO_BIND_ADDRESS="0.0.0.0"
ENTRYPOINT ["/commento/commento"]
