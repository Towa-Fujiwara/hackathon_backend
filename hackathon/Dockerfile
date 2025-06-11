FROM golang:1.24-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ls -la

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./main.go


FROM alpine:latest


WORKDIR /app


COPY --from=builder /app/server .


COPY --from=builder /app/serviceAccountKey.json .


EXPOSE 8080

CMD ["./server"]