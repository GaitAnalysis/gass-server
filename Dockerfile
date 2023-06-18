FROM golang:1.20-alpine
WORKDIR /src
COPY go.mod go.sum ./
COPY . .
RUN go build -o ./server cmd/main.go

FROM alpine
ENV PORT=80
COPY --from=0 /src/server /usr/bin/server
CMD ["server"]