FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=build /bin/app /app
COPY config.toml /
COPY db/migrations/0001_init.sql /init.sql

ARG PG_DSN
ARG JWT_SECRET
ARG TELEGRAM_BOT_TOKEN

ENV PG_DSN=$PG_DSN
ENV JWT_SECRET=$JWT_SECRET
ENV TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN

USER nonroot

EXPOSE 8081
ENTRYPOINT ["/app"]
