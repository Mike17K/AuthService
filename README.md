# Build
```Dockerfile
sudo docker compose up
```

# Localy
Install dependencies
```bash
go mod tidy
```

run with air and .env
```bash
export $(egrep -v ‘^#’ .env | xargs) && air
```
