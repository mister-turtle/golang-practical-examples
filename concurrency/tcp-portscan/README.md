# Portscan Pipeline

This is a quick example to show using channels and waitgroups in a fan-out pattern.

### Usage

```
go run main.go

Usage of <utility>:

  -e int
        Ending port (default 1000)
  -i string
        Target IP address
  -s int
        Starting port (default 1)
  -t int
        Number of threads to run (default 5)
```

### Outline

The program takes an IP address via the flag `-i`  
It can additionally take start and ending ports with `-i` and `-e`  
It will create a struct of each of these and pass them to a channel, which is read by multiple workers.  
The default number of workerse is `5` but can be specified with `-t`  

### Example

```
$ time go run main.go -i 127.0.0.1 -s 1 -e 65535

Simple Concurrent TCP Port Scanner
Starting 5 threads
127.0.0.1        : 9999   - open
127.0.0.1        : 35471  - open
Finished.

real    0m0.422s
user    0m0.553s
sys     0m0.527s
```