FROM golang:1.24-alpine AS deps

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24-alpine AS build

WORKDIR /src

COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY go.mod go.sum ./
COPY back ./back

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/neohome-api ./back/cmd/server

FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata && addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=build /out/neohome-api /app/neohome-api

USER app

EXPOSE 3434

ENTRYPOINT ["/app/neohome-api"]
