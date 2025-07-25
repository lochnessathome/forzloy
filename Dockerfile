FROM dockerhub.timeweb.cloud/golang:1.24.5

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./...

EXPOSE 80

CMD ["app"]

