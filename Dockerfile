FROM golang:1.17-buster AS builder
RUN mkdir -p /build
COPY . /build
RUN cd /build && go build -o plugin.run

FROM debian:bullseye-slim
COPY --from=builder /build/plugin.run /bin/
ENTRYPOINT ["/bin/plugin.run"]
