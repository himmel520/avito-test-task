FROM golang:1.21.6 AS builder

WORKDIR /usr/local/src

COPY go.mod ./ 
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app cmd/app/main.go

FROM ubuntu as runner

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/bin/migrations /

EXPOSE 8080
CMD ["/app"]
