# Cross Routine Channels
  
### Overview
This simple example shows how to use channels to communicate across goroutines.

One function is run in a go routine, taking two integer channels. One channel it uses for input, and one it uses for output.
A second function is started in a go routine that takes the same channel used as an output of the previous goroutine.
  
This allows us to communicate across the go routines by sending data on the shared channel.

### Running

```shell
git clone https://github.com/go-practical-examples/channels.git
cd channels/xroutine
go run xroutine.go
```
