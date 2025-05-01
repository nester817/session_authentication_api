FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod download 
RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /root
ENV DB_URL=""
ENV CACHE_URL=""
COPY --from=builder app/main .
CMD [ "./main" ]