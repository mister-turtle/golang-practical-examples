# Cross Routine Channels
  
### Overview
This simple example shows how to use channels to communicate across goroutines.

One function is run in a go routine, taking two integer channels. One channel it uses for input, and one it uses for output.
A second function is started in a go routine that takes the same channel used as an output of the previous goroutine.
  
This allows us to communicate across the go routines by sending data on the shared channel.

### Running

```shell
go run main.go

Go Channels Demo

 -----Sending new values-----
Above 500 Handler Started
Int Handler Started

 -----Sending new values-----
Int Handler received: 258

 -----Sending new values-----
Int Handler received: 520
Above 500 received: 520

 -----Sending new values-----
Int Handler received: 571
Above 500 received: 571

 -----Sending new values-----
Int Handler received: 315

 -----Sending new values-----
Int Handler received: 56

 -----Sending new values-----
Int Handler received: 258

 -----Sending new values-----
Int Handler received: 142

 -----Sending new values-----
Int Handler received: 280

 -----Sending new values-----
Int Handler received: 503
Above 500 received: 503
Exiting gracefully.
Int Handler received: 304
```
