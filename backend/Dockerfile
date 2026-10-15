# Etapa 1: compilar el binario de la API Go/Fiber.
FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api .

# Etapa 2: imagen mínima de ejecución.
FROM alpine:3.22
WORKDIR /app
COPY --from=build /api /app/api
EXPOSE 3000
CMD ["/app/api"]
