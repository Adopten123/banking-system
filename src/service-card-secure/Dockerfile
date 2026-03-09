FROM rust:alpine AS builder

RUN apk add --no-cache musl-dev protoc protobuf-dev

WORKDIR /app

COPY . .

RUN cargo build --release

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/target/release/service-card-secure .

# COPY --from=builder /app/config ./config

EXPOSE 50051

CMD ["./service-card-secure"]