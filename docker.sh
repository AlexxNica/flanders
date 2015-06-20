CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/flanders main/main.go
docker build -t flanders-small -f Dockerfile.small .
