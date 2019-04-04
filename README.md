Alerta Notifications

## Prerequisites
- setup an Alerta instance
- setup an Smtp server
- setup a slack subscription

## Build and Run
```
go build
./notifications config/config.yaml
```

## Docker

### Build an image

docker build -t "guanaco/notifications:0.0.1" .

### Run Compose file or run container
```bash
cd docker
docker-compose up
```
or
```bash
docker run -td guanaco/notifications:0.0.1
```


