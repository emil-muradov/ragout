FROM golang:1.25.4-alpine

LABEL maintainer="Emil Muradov <emil.muradov22@gmail.com>"
LABEL version="0.1"

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY . .

RUN go mod download

USER appuser

CMD ["go", "run", "cmd/server.go"]