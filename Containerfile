FROM docker.io/library/golang:1.23-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o healthpose healthpose.go types.go

FROM docker.io/library/alpine:3.20

EXPOSE 8080

RUN apk --no-cache add ca-certificates

COPY misc/config /config
COPY --from=build /app/healthpose /usr/local/bin/healthpose

USER nobody
VOLUME [ "/config" ]
CMD ["/usr/local/bin/healthpose"]
