FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o my-s3-clone

FROM debian:bookworm-slim

COPY --from=builder /app/my-s3-clone /my-s3-clone

EXPOSE 8080

ENV MINIO_ENDPOINT=http://minio:9000
ENV MINIO_ACCESS_KEY=minioadmin
ENV MINIO_SECRET_KEY=minioadmin

CMD ["/my-s3-clone"]
