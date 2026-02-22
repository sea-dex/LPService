FROM golang:1.23-alpine as builder

WORKDIR /build

ARG GOPROXY
ENV GOPROXY=${GOPROXY:-} GOCACHE=/root/.cache/go-build CGO_ENABLED=1

COPY go.mod .
COPY go.sum .
# RUN apk add --no-cache gcc musl-dev
RUN apk add git make pkgconf musl-dev gcc librdkafka-dev && go mod download -x
COPY . .

RUN echo "Git commit hash: "$(git rev-parse HEAD)
RUN make

# app
FROM golang:1.23-alpine

RUN apk --no-cache add curl

WORKDIR /app
COPY --from=builder /build/lpservice ./
ENV TZ=UTC
