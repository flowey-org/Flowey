FROM docker.io/node:22.8.0 AS builder

WORKDIR /build/src

COPY . .

RUN --mount=type=cache,target=/root/.npm \
	npm install && \
	npm run build -- --outDir /build/output

FROM docker.io/nginx:1.27.2-alpine-slim

COPY --from=builder /build/output /usr/share/nginx/html
