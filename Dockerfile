FROM golang:latest

WORKDIR /server

COPY . .

RUN go get

RUN chmod +x ./scripts/build.sh

RUN ./scripts/build.sh

CMD ["./build/server"]