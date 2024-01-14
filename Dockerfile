FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download && go build -o goweb .

EXPOSE 3000

CMD ["./goweb"]
