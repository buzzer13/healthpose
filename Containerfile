FROM golang:1.21-alpine3.19 AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux go build -o healthpose healthpose.go types.go

FROM alpine:3.19

EXPOSE 8080

RUN apk --no-cache add ca-certificates libcap

COPY misc/config /config
COPY --from=build /app/healthpose /usr/local/bin/healthpose

RUN setcap cap_net_raw=+ep /usr/local/bin/healthpose

USER nobody
VOLUME [ "/config", "/data" ]
CMD ["/usr/local/bin/healthpose"]
