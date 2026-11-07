FROM casjaysdev/go:latest AS builder

ARG TARGETARCH
ARG VERSION=dev
ARG BUILD_DATE
ARG BUILD_EPOCH
ARG COMMIT_ID
ARG OFFICIAL_SITE

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -buildvcs=false -trimpath \
    -ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.CommitID=${COMMIT_ID}' -X 'main.BuildDate=${BUILD_DATE}'" \
    -o /app/binary/zipcodes ./src

FROM alpine:latest

ARG VERSION=dev
ARG BUILD_DATE
ARG BUILD_EPOCH
ARG COMMIT_ID
ARG LICENSE=MIT

RUN apk add --no-cache git curl bash tini tor ca-certificates tzdata

COPY --from=builder /app/binary/zipcodes /usr/local/bin/zipcodes
COPY docker/rootfs/ /
COPY docker/Dockerfile /root/Dockerfile

RUN chmod 755 /usr/local/bin/*

ENV TZ=America/New_York \
    DATABASE_DIR=/data/db/sqlite

EXPOSE 80
STOPSIGNAL SIGRTMIN+3

HEALTHCHECK --start-period=10m --interval=5m --timeout=15s --retries=3 \
    CMD /usr/local/bin/zipcodes --status || exit 1

ENTRYPOINT [ "tini", "-p", "SIGTERM", "--", "/usr/local/bin/entrypoint.sh" ]
