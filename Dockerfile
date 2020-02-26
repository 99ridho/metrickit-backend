# Builder
FROM golang:1.13.4 as builder
WORKDIR /app
COPY . .
RUN go get -u -v && go mod tidy
RUN CGO_ENABLED=0 go build -o mb .

# Distribution
FROM alpine:latest
WORKDIR /app
RUN apk update && apk upgrade && \
    apk --no-cache --update add ca-certificates tzdata
COPY --from=builder /app/mb .
COPY --from=builder /app/config* ./
RUN chmod +x ./mb
EXPOSE 8185
CMD ["/app/mb"]