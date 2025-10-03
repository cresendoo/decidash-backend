ARG GOLANG_VERSION=1.24

FROM golang:${GOLANG_VERSION}-alpine AS builder

RUN apk update && apk add --no-cache alpine-sdk git 

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY cmd cmd
COPY config config
COPY internal internal
COPY pkg pkg

RUN GOOS=linux GOARCH=amd64 go build -o ./api-server ./cmd/api-server

ARG RELEASE
ENV RELEASE=${RELEASE}

ARG GIT_SHA1
ENV GIT_SHA1=${GIT_SHA1}

ARG GIT_REVISION
ENV GIT_REVISION=${GIT_REVISION}

FROM alpine
WORKDIR /app

RUN apk update
RUN apk add --no-cache tzdata curl

RUN cp /usr/share/zoneinfo/UTC /etc/localtime
RUN echo "UTC" > /etc/timezone

COPY --from=builder /app/api-server /app/api-server

COPY ./entrypoint.sh /app/
ENTRYPOINT ["sh", "/app/entrypoint.sh"]
