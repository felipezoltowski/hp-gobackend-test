# Etapa de build
FROM golang:1.23-bookworm as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o webserver
EXPOSE 8080
CMD ["./webserver"]

# Etapa de execução
FROM golang:1.23-bookworm 
WORKDIR /app
COPY --from=builder /app/webserver .
EXPOSE 8080
CMD ["./webserver"]