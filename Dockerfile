FROM golang:latest AS builder

WORKDIR /app

RUN export GO111MODULE=on

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o ./api ./cmd/api

FROM alpine:latest AS runner

COPY --from=builder /app/api .

EXPOSE 9000 5432 9091 9092 9093

CMD ["./api" ]
