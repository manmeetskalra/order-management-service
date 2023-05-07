FROM golang:1.16

WORKDIR /app

COPY . .

RUN go build -o order-management .

EXPOSE 8080

CMD ["./order-management"]
