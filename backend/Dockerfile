FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
RUN go get github.com/cortesi/modd/cmd/modd
CMD ["modd"]

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
