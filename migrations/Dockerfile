FROM gomicro/goose AS goose


FROM alpine:latest
COPY --from=goose /bin/goose /bin/goose

WORKDIR /migrations
COPY . /migrations