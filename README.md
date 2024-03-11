# echo-boilerplate [![Go Report Card](https://goreportcard.com/badge/github.com/alexferl/echo-boilerplate)](https://goreportcard.com/report/github.com/alexferl/echo-boilerplate) [![codecov](https://codecov.io/gh/alexferl/echo-boilerplate/branch/master/graph/badge.svg)](https://codecov.io/gh/alexferl/echo-boilerplate)

A Go 1.22+ boilerplate app using the minimalist [echo](https://github.com/labstack/echo) framework and with
authentication, authorization and request/response validation.

> **Note**: I use this as a starting point for personal projects, it can and will change without notice.

## Features
- [JWT](https://jwt.io/) for authentication with access and [refresh](https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/) tokens.
 The access token can be sent in the [Authorization](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization) header or
 as a [cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies). See [echo-jwt](https://github.com/alexferl/echo-jwt).
- [Casbin](https://casbin.io/) for authorization using RBAC. See [echo-casbin](https://github.com/alexferl/echo-casbin).
- [OpenAPI](https://www.openapis.org/) for request and response validation. See [echo-openapi](https://github.com/alexferl/echo-openapi).

## Requirements
Before getting started, install the following:

Required:
- [pre-commit](https://pre-commit.com/#install)
- [MongoDB](https://www.mongodb.com/docs/manual/installation/#mongodb-installation-tutorials)

Optional:

- [gofumpt](https://pkg.go.dev/mvdan.cc/gofumpt) (needed to run `make fmt`)
- [redocly-cli](https://redocly.com/docs/cli/installation/) (needed to run `make openapi-lint`)

## Using
Setup the dev environment first:
```shell
make dev
```
>**Note**: An RSA private key will be generated in the current folder to sign and verify the JSON web tokens.

### Creating the superuser
Launch the superuser cmd with `go run ./cmd/superuser --password <your password>`. You can change the default values
with the following flags: `--email`, `--name` and `--username`. You can view all the other settings with `--help`.

### Building & Running locally
```shell
make run
```
### Using the API
#### Login
Request:
```shell
curl --request POST \
  --url http://localhost:1323/auth/login \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "super@example.com",
	"password": "<your password>"
}'
```
Response:
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer"
}
```
**Note**: The `access_token` only lasts 60 minutes by default, this is as designed. A client
(like an [SPA](https://en.wikipedia.org/wiki/Single-page_application) or a mobile application) would have an interceptor
to catch the 401 responses, send the `refresh_token` to the `/auth/refresh` endpoint to get new access and refresh tokens and
then retry the previous request with the new `access_token` which should then succeed. The duration of the `access_token`
can be modified with `--jwt-access-token-expiry` and the `refresh_token` with `--jwt-refresh-token-expiry`.

#### Get currently authenticated user
Request:

Using the `Authorization` header:
```shell
curl --request GET \
  --url http://localhost:1323/me \
  --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...'
```

Using the cookie (the cookie is sent automatically with web browsers, HTTPie and some other clients):
```shell
curl --request GET \
  --url http://localhost:1323/me \
  --cookie access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

Response:
```json
{
  "id": "cdhgh0dfclscplnrcuag",
  "username": "super",
  "email": "super@example.com",
  "name": "Super",
  "bio": "",
  "created_at": "2022-11-03T00:17:05.837Z",
  "updated_at": null
}
```

### OpenAPI docs
You can see the OpenAPI docs by running the app and navigating to `http://localhost:1323/docs` or by
opening [assets/index.html](docs/index.html) in your web browser.

### Repository layout
```
.
├── casbin    <--- model and policy files for Casbin
├── cmd       <--- entrypoints
├── config    <--- config structs and defaults are specified here
├── configs   <--- config files, for configs that rarely change, but should override the defaults
├── data      <--- base mapper, database helpers
├── docs      <--- generated documentation from OpenAPI schema
├── handlers  <--- HTTP handlers (aka controllers, endpoints etc.) that interacts with the services
├── mappers   <--- mapper layer that the services use to insert/retrieve models from the database
├── models    <--- structs defining the various resources
├── openapi   <--- OpenAPI schema files
├── server.go <--- glues handlers/services/mappers
├── services  <--- service layer that interacts with the mappers
├── testing   <--- testing helpers
└── util      <--- general helpers
```

### Usage
```shell
go run ./cmd/server --help
```

```
Usage of ./echo-boilerplate:
      --app-name string                                The name of the application. (default "app")
      --base-url string                                Base URL where the app will be served (default "http://localhost:1323")
      --casbin-model string                            Casbin model file (default "./casbin/model.conf")
      --casbin-policy string                           Casbin policy file (default "./casbin/policy.csv")
      --cookies-domain string                          Cookies domain
      --cookies-enabled                                Send cookies with authentication requests
      --csrf-cookie-domain string                      CSRF cookie domain
      --csrf-cookie-name string                        CSRF cookie name (default "csrf_token")
      --csrf-enabled                                   CSRF enabled
      --csrf-header-name string                        CSRF header name (default "X-CSRF-Token")
      --csrf-secret-key string                         CSRF secret used to hash the token
      --env-name string                                The environment of the application. Used to load the right configs file. (default "local")
      --http-bind-address ip                           The IP address to listen at. (default 127.0.0.1)
      --http-bind-port uint                            The port to listen at. (default 1323)
      --http-cors-allow-credentials                    Tells browsers whether to expose the response to frontend JavaScript code when the request's credentials mode (Request.credentials) is 'include'.
      --http-cors-allow-headers strings                Indicate which HTTP headers can be used during an actual request.
      --http-cors-allow-methods strings                Indicates which HTTP methods are allowed for cross-origin requests. (default [GET,HEAD,PUT,PATCH,POST,DELETE])
      --http-cors-allow-origins strings                Indicates whether the response can be shared with requesting code from the given origin. (default [*])
      --http-cors-enabled                              Enable cross-origin resource sharing.
      --http-cors-expose-headers strings               Indicates which headers can be exposed as part of the response by listing their name.
      --http-cors-max-age int                          Indicates how long the results of a preflight request can be cached.
      --http-graceful-timeout duration                 Timeout for graceful shutdown. (default 30s)
      --http-log-requests                              Controls the logging of HTTP requests (default true)
      --http-tls-cert-file string                      TLS certificate file
      --http-tls-key-file string                       TLS key file
      --jwt-access-token-cookie-name string            JWT access token cookie name (default "access_token")
      --jwt-access-token-expiry duration               JWT access token expiry (default 1h0m0s)
      --jwt-issuer string                              JWT issuer (default "http://localhost:1323")
      --jwt-private-key string                         JWT private key file path (default "./private-key.pem")
      --jwt-refresh-token-cookie-name string           JWT refresh token cookie name (default "refresh_token")
      --jwt-refresh-token-expiry duration              JWT refresh token expiry (default 720h0m0s)
      --log-level string                               The granularity of log outputs. Valid levels: 'PANIC', 'FATAL', 'ERROR', 'WARN', 'INFO', 'DEBUG', 'TRACE', 'DISABLED' (default "INFO")
      --log-output string                              The output to write to. 'stdout' means log to stdout, 'stderr' means log to stderr. (default "stdout")
      --log-writer string                              The log writer. Valid writers are: 'console' and 'json'. (default "console")
      --mongodb-app-name string                        MongoDB app name
      --mongodb-connect-timeout-ms duration            MongoDB connect timeout ms (default 10s)
      --mongodb-password string                        MongoDB password
      --mongodb-replica-set string                     MongoDB replica set
      --mongodb-server-selection-timeout-ms duration   MongoDB server selection timeout ms (default 10s)
      --mongodb-socket-timeout-ms duration             MongoDB socket timeout ms (default 30s)
      --mongodb-uri string                             MongoDB URI (default "mongodb://localhost:27017")
      --mongodb-username string                        MongoDB username
      --oauth2-google-client-id string                 OAuth2 Google client id
      --oauth2-google-client-secret string             OAuth2 Google client secret
      --oauth2-providers strings                       OAuth2 providers
      --openapi-schema string                          OpenAPI schema file (default "./openapi/openapi.yaml")
```

### Docker
#### Build
```shell
make docker-build
```

#### Run
```shell
make docker-run
```

#### Passing args
CLI:
```shell
docker run -p 1323:1323 --rm echo-boilerplate --env-name prod
```

Environment variables:
```shell
docker run -p 1323:1323 -e "APP_ENV_NAME=prod" --rm echo-boilerplate
```
