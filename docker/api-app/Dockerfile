FROM golang:latest

LABEL maintainer="Dmitrii Kudryavtsev <dmitrii@ordercube.de>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o ./bin/api-app ./cmd/api/main.go

#RUN go build -o ./bin/api-app ./cmd/api/main.go

CMD ["./bin/api-app"]
