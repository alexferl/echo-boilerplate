# echo-boilerplate [![Go Report Card](https://goreportcard.com/badge/github.com/admiralobvious/echo-boilerplate)](https://goreportcard.com/report/github.com/admiralobvious/echo-boilerplate)

A small boilerplate app using the minimalist [echo](https://github.com/labstack/echo) framework
with [12-factor](https://12factor.net/) and following golang-standards' [project-layout](https://github.com/golang-standards/project-layout).

### Building & Running locally
```shell script
go build ./cmd/app && ./app
```

### Usage
```shell script
Usage of ./app:
      --app-name string               The name of the application. Used to prefix environment variables. (default "app")
      --bind-address ip               The IP address to listen at. (default 127.0.0.1)
      --bind-port uint                The port to listen at. (default 1323)
      --cors-allow-credentials        Tells browsers whether to expose the response to frontend JavaScript code when the request's credentials mode (Request.credentials) is 'include'.
      --cors-allow-headers strings    Indicate which HTTP headers can be used during an actual request.
      --cors-allow-methods strings    Indicates which HTTP methods are allowed for cross-origin requests. (default [GET,HEAD,PUT,PATCH,POST,DELETE])
      --cors-allow-origins strings    Indicates whether the response can be shared with requesting code from the given origin. (default [*])
      --cors-expose-headers strings   Indicates which headers can be exposed as part of the response by listing their name.
      --cors-max-age int              Indicates how long the results of a preflight request can be cached.
      --env-name string               The environment of the application. Used to load the right config file. (default "local")
      --graceful-timeout uint         Timeout for graceful shutdown. (default 30)
      --log-file string               The log file to write to. 'stdout' means log to stdout, 'stderr' means log to stderr and 'null' means discard log messages. (default "stdout")
      --log-format string             The log format. Valid format values are: text, json. (default "text")
      --log-level string              The granularity of log outputs. Valid log levels: debug, info, warning, error and critical. (default "info")
      --log-requests-disabled         Disables HTTP requests logging.
```

### Building Docker image
```shell script
docker build -t app .
```
