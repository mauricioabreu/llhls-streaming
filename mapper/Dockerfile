FROM golang:1.23-alpine AS app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mapper .

FROM alpine:3.20

WORKDIR /root

COPY --from=app /app/mapper .

CMD ["./mapper"]
