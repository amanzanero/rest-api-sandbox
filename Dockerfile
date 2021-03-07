FROM golang:1.15-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build


FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /app/rest-api-sandbox /app/
EXPOSE 8080
ENTRYPOINT ["/app/rest-api-sandbox"]
