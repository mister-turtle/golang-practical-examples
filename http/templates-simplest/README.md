# Simplest-Templates

This example shows the basic usage of templates for an http service.
A route is established using the stdlib net/http router which renders the clients IP address and current server time into the response.  
We use a closure to access the template from inside the handler function, and a map to pass the data into the template.  
We then make use of the `io.Writer` interface and pass the `http.ResponseWriter` directly to the template function to write the result.  

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
2021/06/29 22:15:03 Listening on 127.0.0.1:8000
```
```
$ curl 127.0.0.1:8000

<html>
<head>
        <title>Welcome, 127.0.0.1:48130</title>
</head>
<body>
        <h1>Hello, 127.0.0.1:48130</h1>
        <br>
        <p>The current server time is 2021-06-29 22:15:23.4158511</p>
</body>
</html>
```
