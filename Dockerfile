FROM golang:1.22-bullseye

# Development image
WORKDIR /app
ENV CGO_ENABLED=0
RUN go install github.com/cosmtrek/air@v1.49.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.22.1
COPY go.mod go.sum ./

RUN go mod download
CMD ["air"]
