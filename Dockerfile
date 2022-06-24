# Build stage
FROM golang:alpine3.15 AS builder
RUN apk add git
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

COPY . .
RUN go build -o deployment

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/deployment .

EXPOSE 8080
CMD ["/app/deployment"]
