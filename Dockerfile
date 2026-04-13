FROM golang:1.26-alpine AS builder

WORKDIR /usr/src/api

COPY go.* ./

RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -trimpath --ldflags="-s -w" -v -o /usr/local/bin/api ./cmd/api/main.go

# Using a distroless image to reduce the image size and improve security.
# If we use FROM scratch here, the image will not have any ca-certificates and therefore https calls to external services will fail. A distroless image contains ca-certificates so we can use it.
# See: https://github.com/GoogleContainerTools/distroless
# Also: https://github.com/GoogleContainerTools/distroless/tree/main/base

FROM gcr.io/distroless/static-debian12:latest-amd64

COPY --from=builder /usr/local/bin/api /usr/local/bin/api

COPY --from=builder /usr/src/api/migrations /migrations

CMD ["api"]
