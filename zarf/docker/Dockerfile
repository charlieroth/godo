FROM golang:1.23 as builder
ENV CGO_ENABLED=0 
ARG BUILD_REF

COPY . /app

WORKDIR /app/cmd/server
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go binary in Alpine
FROM alpine:3.21
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S zarf && \
    adduser -u 1000 -h /app -G zarf -S zarf
COPY --from=builder --chown=zarf:zarf /app/cmd/server/server /app/server

WORKDIR /app
USER zarf
CMD ["./server"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
    org.opencontainers.image.title="godo" \
    org.opencontainers.image.authors="Charlie Roth <charlieroth4@gmail.com>" \
    org.opencontainers.image.source="https://github.com/charlieroth/godo" \
    org.opencontainers.image.revision="${BUILD_REF}" \
    org.opencontainers.image.vendor="Charlie Roth"
