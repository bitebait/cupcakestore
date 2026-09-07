# syntax=docker/dockerfile:1
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /out/cupcakestore .

FROM debian:bookworm-slim AS runtime
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 10001 cupcake \
    && useradd --uid 10001 --gid cupcake --no-create-home --shell /usr/sbin/nologin cupcake
WORKDIR /app
COPY --from=build /out/cupcakestore /usr/local/bin/cupcakestore
COPY views ./views
COPY web ./web
RUN mkdir -p /data /app/web/images \
    && chown -R cupcake:cupcake /data /app/web/images
ENV APP_HOST=0.0.0.0 APP_PORT=8080 DB_TYPE=sqlite DB_PATH=/data/store.db
USER cupcake:cupcake
EXPOSE 8080
STOPSIGNAL SIGTERM
ENTRYPOINT ["cupcakestore"]
