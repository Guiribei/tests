FROM golang:1.27.4-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/caderno ./cmd/server

FROM alpine:3.23
RUN apk add --no-cache ca-certificates git
WORKDIR /app
COPY --from=build /out/caderno /usr/local/bin/caderno
COPY option-one.json option-two.json ./
RUN mkdir -p /data && chown 65532:65532 /data
ENV ADDR=:8080 \
    DATABASE_PATH=/data/provas.db \
    EXAMS_DIR=/app
VOLUME ["/data"]
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/caderno"]

