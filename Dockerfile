FROM golang:1.21-bullseye

ARG USER_ID
ARG GROUP_ID

VOLUME /go /src
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN go build -ldflags='-s -w' ./cmd/rudolphe && \
    chown ${USER_ID}:${GROUP_ID} rudolphe
