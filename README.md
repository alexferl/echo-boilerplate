# echo-boilerplate
A small boilerplate app using the minimalist [echo](https://github.com/labstack/echo) framework
with [12-factor](https://12factor.net/) and following golang-standards' [project-layout](https://github.com/golang-standards/project-layout).

### Building & Running
```shell script
go build ./cmd/app && ./app
```

### Building Docker image
```shell script
docker build -t app .
```
