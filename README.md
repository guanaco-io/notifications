Alerta Notifications

## Prerequisites
- setup an Alerta instance
- setup an Smtp server
- setup a slack subscription

## Build and Run
to run the project locally, adjust `config/config.yml` to your local setup and execute
```
go build
./notifications config/config.yml
```

## Docker

### Build an image
```bash
docker build -t "guanaco/notifications:0.0.1" .
```

### Run Compose file or run container
```bash
cd docker
docker-compose up
```
or
```bash
docker run -td guanaco/notifications:0.0.1
```


