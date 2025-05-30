FROM golang:1.20 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go


FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/app /app/app
COPY --from=builder /app/static /static


USER nonroot:nonroot

ENV PORT=8080
EXPOSE 8080

CMD ["/app/app"]
