# Simplest HTTP Server

This example uses the standard library to host a simple HTTP server. The listen address and port can be changed via `-p` and `-l`. There are a few routes designed to show using `http.ResponseWriter` including setting cookies, headers, and using redirects.

### Usage

```
go run main.go

Usage of <utility>:

  -l string
        IP address to listen on (default "127.0.0.1")
  -p int
        TCP port to listen on (default 8000)
```

### Example
```
$ go run main.go
2021/06/29 21:37:50 Listening on 127.0.0.1:8000
```
```
$ curl http://127.0.0.1:8000/
Hello, world!

$ curl -v http://127.0.0.1:8000/redirect
[...]]
< HTTP/1.1 308 Permanent Redirect
< Content-Type: text/html; charset=utf-8
< Location: /
< Date: Tue, 29 Jun 2021 20:38:50 GMT
< Content-Length: 37
<
<a href="/">Permanent Redirect</a>.
```
```
$ curl -v http://127.0.0.1:8000/headers
[...]]
< HTTP/1.1 200 OK
< Custom-Header: super useful value
< Date: Tue, 29 Jun 2021 20:39:22 GMT
< Content-Length: 18
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 127.0.0.1 left intact
Check the headers!
```
```
$ curl -v http://127.0.0.1:8000/cookie
[...]
< HTTP/1.1 200 OK
< Set-Cookie: MyCookie=123456789
< Date: Tue, 29 Jun 2021 20:39:46 GMT
< Content-Length: 33
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 127.0.0.1 left intact
Check the cookie in the response!
```

