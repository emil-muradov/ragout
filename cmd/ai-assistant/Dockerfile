FROM golang:1.25.1-alpine

MAINTAINER Emil Muradov <emil.muradov22@gmail.com>

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

WORKDIR /app

COPY --chown=appuser:appgroup . .

RUN go mod download

WORKDIR /app/cmd/ai-assistant

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]