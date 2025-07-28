FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.mod
COPY cmd/main.go cmd/main.go

ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o counter cmd/main.go

FROM gcr.io/distroless/static:nonroot AS production

COPY --from=builder /app/counter /counter

USER 65532

CMD ["/counter"] 