FROM golang:1.22-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY ./hackathon ./hackathon


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./hackathon/main.go


FROM alpine:latest


WORKDIR /app


COPY --from=builder /app/server .


COPY ./hackathon/serviceAccountKey.json ./hackathon/serviceAccountKey.json


EXPOSE 8080

CMD ["./server"]