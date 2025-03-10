FROM golang:latest AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/shortURL/bin/main ./cmd/shortURL/

FROM alpine:latest
WORKDIR /shortURL
RUN mkdir /shortURL/logs
COPY --from=build /build/cmd/shortURL/bin/main .
CMD ["/shortURL/main"]