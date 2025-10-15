FROM golang:1.25.1-alpine

LABEL maintainer="Emil Muradov <emil.muradov22@gmail.com>"
LABEL version="0.1"

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

WORKDIR /app

COPY --chown=appuser:appgroup . .

RUN go mod download

WORKDIR /app/cmd/ai-assistant

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]