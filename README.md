# echo-boilerplate [![Go Report Card](https://goreportcard.com/badge/github.com/alexferl/echo-boilerplate)](https://goreportcard.com/report/github.com/alexferl/echo-boilerplate)

A small boilerplate app using the minimalist [echo](https://github.com/labstack/echo) framework
with [12-factor](https://12factor.net/).

### Building & Running locally
```shell
make run
```

### Usage
```shell
Usage of ./app:
      --app-name string                    The name of the application. (default "app")
      --env-name string                    The environment of the application. Used to load the right configs file. (default "local")
      --http-bind-address ip               The IP address to listen at. (default 127.0.0.1)
      --http-bind-port uint                The port to listen at. (default 1323)
      --http-cors-allow-credentials        Tells browsers whether to expose the response to frontend JavaScript code when the request's credentials mode (Request.credentials) is 'include'.
      --http-cors-allow-headers strings    Indicate which HTTP headers can be used during an actual request.
      --http-cors-allow-methods strings    Indicates which HTTP methods are allowed for cross-origin requests. (default [GET,HEAD,PUT,PATCH,POST,DELETE])
      --http-cors-allow-origins strings    Indicates whether the response can be shared with requesting code from the given origin. (default [*])
      --http-cors-enabled                  Enable cross-origin resource sharing.
      --http-cors-expose-headers strings   Indicates which headers can be exposed as part of the response by listing their name.
      --http-cors-max-age int              Indicates how long the results of a preflight request can be cached.
      --http-graceful-timeout uint         Timeout for graceful shutdown. (default 30)
      --http-log-requests-disabled         Disable the logging of HTTP requests
      --log-level string                   The granularity of log outputs. Valid log levels: 'panic', 'fatal', 'error', 'warn', 'info', 'debug' and 'trace'. (default "info")
      --log-output string                  The output to write to. 'stdout' means log to stdout, 'stderr' means log to stderr. (default "stdout")
      --log-writer string                  The log writer. Valid writers are: 'console' and 'json'. (default "console")
```

### Building Docker image
```shell
docker build -t app .
```
