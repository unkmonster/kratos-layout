FROM golang:1.24.6 AS builder

ARG SERVICE_NAME

WORKDIR /src
# cache go mod download
COPY go.mod go.sum .
RUN go mod download

RUN mkdir -p pkg

COPY ./app/${SERVICE_NAME} ./app/${SERVICE_NAME}
COPY ./api ./api
COPY ./.git ./.git
COPY ./third_party ./third_party
COPY ./pkg ./pkg
RUN make -C app/${SERVICE_NAME} build

# ensure the migrations and casbin directory must be exist
RUN mkdir -p ./app/${SERVICE_NAME}/migrations
RUN mkdir -p ./app/${SERVICE_NAME}/casbin

FROM debian:stable-slim

ARG SERVICE_NAME
ARG CMD_NAME

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/app/${SERVICE_NAME}/bin/${CMD_NAME} /app/server
COPY --from=builder /src/app/${SERVICE_NAME}/migrations /data/migrations
COPY --from=builder /src/app/${SERVICE_NAME}/casbin /data/casbin

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

ENV SERVICE_NAME=${SERVICE_NAME}
ENTRYPOINT ["./server"]
CMD ["-conf", "/data/conf"]
