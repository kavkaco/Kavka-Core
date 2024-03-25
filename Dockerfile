FROM golang:1.22-alpine

WORKDIR /server

COPY . .

RUN chmod +x ./scripts/run_prod.sh

CMD ["./scripts/run_prod.sh"]
