FROM golang:1.23.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

FROM alpine:3.20.3

WORKDIR /app

COPY --from=build /api .

EXPOSE 5000

CMD ["./api"]
