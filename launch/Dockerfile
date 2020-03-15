FROM widget-proxy-src as builder
ADD model.go /build/generated/model.go
RUN GOOS=linux GOARCH=386 go build -o main .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]