FROM docker.io/golang:1.23.1-alpine AS builder

WORKDIR /build/src

RUN apk add --no-cache \
	gcc=13.2.1_git20240309-r0 \
	musl-dev=1.2.5-r0

RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=cache,target="/root/.cache/go-build" \
	--mount=type=bind,target=. \
	mkdir -p /build/output && \
	CGO_ENABLED=1 go build -v -ldflags "-s -w" -trimpath -o /build/flowey

FROM alpine:3.20

LABEL org.opencontainers.image.source=https://github.com/flowey-org/Flowey

EXPOSE 80

ARG USER=default
ENV HOME=/home/$USER

RUN adduser -D $USER

COPY --from=builder /build/flowey /usr/local/bin

USER $USER
WORKDIR $HOME

ENTRYPOINT [ "flowey" ]
CMD [ "server" ]
