FROM golang:alpine as builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY ./04/* .
RUN go build -o /tmp/gophercon-talk .

FROM alpine
COPY --from=builder /tmp/gophercon-talk /usr/bin/gophercon-talk
WORKDIR /app
EXPOSE 23234
VOLUME [ "/app/.ssh" ]
ENTRYPOINT [ "/usr/bin/gophercon-talk" ]

