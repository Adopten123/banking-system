FROM rust:alpine AS builder

RUN apk add --no-cache musl-dev protoc protobuf-dev

WORKDIR /app

COPY proto/ ./proto/

COPY src/service-card-secure/ ./src/service-card-secure/

WORKDIR /app/src/service-card-secure

RUN cargo build --release

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/src/service-card-secure/target/release/service-card-secure .
COPY --from=builder /app/src/service-card-secure/config ./config

EXPOSE 50051

CMD ["./service-card-secure"]