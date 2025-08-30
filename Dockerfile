FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
LABEL org.opencontainers.image.source="https://github.com/FIAP-SOAT-G20/hackathon-video-service" \
      org.opencontainers.image.authors="FIAP 10SOAT G19" \
      org.opencontainers.image.title="Fast Food FIAP TC-4" \
      org.opencontainers.image.description="Image of a backend Order API for a fast food restaurant"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -ldflags "-w -s" -o api cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
