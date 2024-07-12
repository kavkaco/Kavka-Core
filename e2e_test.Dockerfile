FROM golang:1.22.4

WORKDIR /server

COPY . .

RUN go mod tidy

ENV KAVKA_END=test

CMD ["make", "e2e_test"]
