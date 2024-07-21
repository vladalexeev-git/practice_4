FROM golang:1.22 as builder
WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build -o city-linux .
RUN chmod +x city-linux

FROM ubuntu:22.04
# Установка необходимых инструментов и зависимостей
RUN apt-get update && \
        apt-get install -y \
        gcc \
        musl-tools \
        libaio1 \
        && apt-get clean && \
        rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/city-linux .
RUN chmod +x city-linux

CMD ["/app/city-linux"]

