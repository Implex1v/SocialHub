FROM golang:1.19-bullseye as build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go test -coverage ./...
RUN go build -o /app/youtube-worker

FROM debian:bullseye-slim
WORKDIR /app

ARG BUILD_VERSION=0.1.0
ARG GIT_COMMIT=unknown

LABEL org.opencontainers.image.version="${BUILD_VERSION}"\
      org.opencontainers.image.revision="${GIT_COMMIT}"

ENV APP_VERSION ${BUILD_VERSION}

COPY --from=build /app/youtube-worker /app/youtube-worker
EXPOSE 8000

CMD [ "/app/youtube-worker" ]