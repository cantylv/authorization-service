# Builder state
FROM golang:1.23.2 AS builder
RUN apt-get update && apt-get install -y make git curl && apt-get clean

ARG MODULE_NAME=sbertech_backend
WORKDIR /home/${MODULE_NAME}

COPY . .

# building exe ile
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main/main.go

# Production state
FROM alpine:3.20.3 as production
WORKDIR /root/
ARG MODULE_NAME_BUILDER=sbertech_backend

COPY --from=builder /home/${MODULE_NAME_BUILDER}/config/config.yaml config/config.yaml
COPY --from=builder /home/${MODULE_NAME_BUILDER}/main .

RUN chown root:root main

CMD ["./main"]