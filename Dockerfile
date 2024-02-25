# syntax=docker/dockerfile:1
FROM golang:latest
WORKDIR /app
ADD . .
RUN go build -o ./out/go-practise-project ./cmd/food-delivery-main
EXPOSE 8000
ENTRYPOINT ["./out/go-practise-project"]