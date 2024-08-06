######## Builder Base #######
FROM golang:1.22-alpine AS builder-base

RUN apk add --no-cache upx \
  && mkdir /build

WORKDIR /build

######## Builder 1 #######
FROM builder-base AS builder1

COPY client/ /build/

RUN go mod download \
  # -s and -w linker flags for striping debugging information
  && GOOS=linux go build -ldflags="-s -w" -o client . \
  # Addititional UPX compression at the expense of a 10x startup time
  && upx client

######## Builder 2 #######
FROM builder-base AS builder2

COPY api/ /build/

RUN go mod download \
  && GOOS=linux go build -ldflags="-s -w" -o api . \
  && upx api

######## Builder 3 #######
FROM builder-base AS builder3

COPY montage/ /build/

RUN go mod download \
  # https://github.com/gographics/imagick?tab=readme-ov-file#build-tags
  && GOOS=linux go build -tags no_pkgconfig -ldflags="-s -w" -o montage . \
  && upx montage

######## Start a new stage from scratch #######
FROM alpine:latest

ARG UID="${UID:-1000}"
ARG GID="${GID:-1000}"
ARG UNAME="${UNAME:-mbrav}"
ARG GNAME="${GNAME:-mbrav}"

RUN apk add --upgrade --latest --no-cache \
  ca-certificates \
  imagemagick-dev \
  && apk clean cache \
  && addgroup -S -g "${GID}" "${GNAME}" \
  && adduser -S -G "${GNAME}" -u "${UID}" "${UNAME}" -s /bin/ash

COPY --from=builder1 /build/client /usr/local/bin/
COPY --from=builder2 /build/api /usr/local/bin/
COPY --from=builder3 /build/montage /usr/local/bin/

WORKDIR /usr/local/bin/
USER "${UNAME}"

