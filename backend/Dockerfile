FROM golang:1.23.1 AS builder

WORKDIR /app/src

RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=cache,target="/root/.cache/go-build" \
	--mount=type=bind,target=. \
	CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/flowey

FROM scratch

COPY --from=builder /app/flowey /app/flowey

ENTRYPOINT [ "/app/flowey" ]
