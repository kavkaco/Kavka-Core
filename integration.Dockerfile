FROM golang:1.22-alpine

WORKDIR /server

COPY . .

RUN chmod +x ./scripts/integration_test.sh

CMD ["./scripts/integration_test.sh"]
