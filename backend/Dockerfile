FROM docker.io/golang:1.23.1 AS builder

WORKDIR /app/src

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=cache,target="/root/.cache/go-build" \
	CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/flowey

FROM scratch

COPY --from=builder /app/flowey /app/flowey

ENTRYPOINT [ "/app/flowey" ]
