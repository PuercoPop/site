FROM golang:1.19.1-alpine as builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build ./...

FROM golang:1.19.1-alpine
WORKDIR /app
COPY --from=builder /build/site /app
CMD ["site"]
EXPOSE 8080
