# backend build (api server)
FROM golang:1.15-alpine AS api-build
RUN apk add --no-cache --update bash dep make git curl g++ wget ca-certificates

# Configure Go environment for better network reliability
ENV GOPROXY="https://proxy.golang.org,direct"
ENV GOSUMDB=off
ENV GO111MODULE=on
ENV CGO_ENABLED=0

ARG RELEASE=prod
COPY ./api /go/src/commento/api/
WORKDIR /go/src/commento/api

# Pre-download dependencies with retry logic and handle sync module issue
RUN go mod download || go mod download || go mod download

# Skip go mod tidy and directly build
RUN for i in $(seq 1 5); do \
        echo "Build attempt $i..." && \
        make ${RELEASE} -j$(($(nproc) + 1)) && exit 0 || \
        echo "Build failed, retrying in 5 seconds..." && \
        sleep 5; \
    done && exit 1


# frontend build (html, js, css, images)
FROM node:16-alpine AS frontend-build
RUN apk add --no-cache --update bash make g++

ARG RELEASE=prod
COPY ./frontend /commento/frontend
WORKDIR /commento/frontend/
RUN npm install && npm run ${RELEASE}


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
