FROM golang:1.23.1-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

EXPOSE 3000

CMD [ "/api" ]