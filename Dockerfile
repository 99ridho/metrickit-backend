# Builder
FROM golang:1.13.4 as builder
WORKDIR /go/src/github.com/99ridho/metrickit-backend
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure -v && CGO_ENABLED=0 go build -o mb .

# Distribution
FROM alpine:latest
WORKDIR /app
RUN apk update && apk upgrade && \
    apk --no-cache --update add ca-certificates tzdata
COPY --from=builder /go/src/github.com/99ridho/metrickit-backend/mb .
COPY --from=builder /go/src/github.com/99ridho/metrickit-backend/config* ./
RUN chmod +x ./mb
ENV MB_ENV=DOCKER
EXPOSE 8185
CMD ["/app/mb"]