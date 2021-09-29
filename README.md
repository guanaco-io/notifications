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

## Release
Find the latest tag:
```shell
git fetch --tags
git describe --tags $(git rev-list --tags --max-count=1)
```
Create a new tag and publish docker image to Github packages
```shell
./release.sh <release version> <github_username> <github_packages_token>
```

## Docker

### Build an image
```bash
docker build -t "guanacoio/notifications:latest" .
```

### Run Compose file or run container
```bash
cd docker
docker-compose up
```
or
```bash
docker run -td guanacoio/notifications:latest
```


