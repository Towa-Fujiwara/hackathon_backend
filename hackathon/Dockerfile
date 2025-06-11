FROM golang:1.24-alpine AS builder


WORKDIR /app


COPY hackathon/go.mod hackathon/go.sum ./
RUN go mod download

COPY hackathon/ .


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./main.go


FROM alpine:latest


WORKDIR /app


COPY --from=builder /app/server .


COPY --from=builder /app/serviceAccountKey.json .


EXPOSE 8080

CMD ["./server"]