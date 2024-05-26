FROM golang:1.22-alpine

WORKDIR /server

COPY . .

CMD ["make", "build"]
