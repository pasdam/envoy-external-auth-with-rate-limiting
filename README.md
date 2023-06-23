# External auth and rate limiting with Envoy

This repo contains a playground to test external authentication and rate
limiting flow in [Envoy proxy](https://envoyproxy.io/).

There are 2 different versions of the playground one that uses directly
[Envoy](https://envoyproxy.io/) as ingress and another one that uses
[Gloo Edge](https://www.solo.io/products/gloo-edge/).

## Components

The playground components are:

* [Backend](./backend/): it processes the request returning a simple hello
  message;
* [External authentication service](./extauth/): it verifies the authentication
  token;
* [Rate limiter](https://github.com/envoyproxy/ratelimit): it verify the rate of
  the requests using the header added by the auth service; it uses redis to
  cache the status; for this example we use the implementation provided by envoy
  itself;
* [Envoy](https://envoyproxy.io/): it works as ingress and verify the request
  using the external authentication service and rate limiter;
* [Gloo Edge](https://www.solo.io/products/gloo-edge/): it works as API gateway,
  it's basically a wrapper around Envoy.

## Flow

The call flow is:

1. client call endpoint with bearer token in the header;
2. envoy verify the request against the external authentication service: if the
   token in the header is valid it proceed with the following steps, otherwise
   it will return a `401 Unauthorized` error immediately without passing the
   request to the backend;
3. envoy will verify the request with the rate limiter: if the request is
   invalid it will return `` immediately, otherwise will pass the request to
   the backend;
4. the backend process the request and returns a response.

## Usage

### Build and run

To build and run the playground with Envoy just run the following:

```sh
docker compose up
```

To build and run the the one with Gloo:

```sh
docker compose -f docker-compose-gloo.yaml up
```

### Test authentication

To call the backend a bearer token must be passed in, e.g.

```sh
$ curl -v -H "Authorization: Bearer user-1-token" http://localhost:8080
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /slowpath HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
> Authorization: Bearer user-1-token
>
< HTTP/1.1 200 OK
< date: Wed, 20 Jan 2021 18:55:34 GMT
< content-length: 198
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 0
< server: envoy
<
Hello user-1!
X-Envoy-Expected-Rq-Timeout-Ms: 15000
Content-Length: 0
User-Agent: curl/7.64.1
Accept: */*
X-Forwarded-Proto: http
X-Request-Id: 9bd6d869-5a7e-479e-9557-ba0485c56559
* Connection #0 to host localhost left intact
X-User-Id: user-1* Closing connection 0
```

If the token is invalid the response is a `401`.

### Test rate limiting

To test the rate limiting we can use the `/slowpath` endpoint which has a policy
of 2 requests per minute, after that it will reply with:

```sh
$ curl -v -H "Authorization: Bearer user-1-token" http://localhost:8080/slowpath
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8010 (#0)
> GET /slowpath HTTP/1.1
> Host: localhost:8010
> User-Agent: curl/7.64.1
> Accept: */*
> Authorization: Bearer user-1-token
>
< HTTP/1.1 429 Too Many Requests
< x-envoy-ratelimited: true
< date: Wed, 20 Jan 2021 18:55:42 GMT
< server: envoy
< content-length: 0
<
* Connection #0 to host localhost left intact
* Closing connection 0
```
