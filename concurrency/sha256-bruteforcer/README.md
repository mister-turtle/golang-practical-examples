# SHA256 Bruteforcer

This is a quick example to show using channels, waitgroups, and context cancellation in a "real world utility".  
The intention isn't to be a *useful* brute forcing utility, although suprisingly it managed 27 million h/s on my laptop.

### Usage

```
go run main.go

Usage of <utility>:

  -h string
        SHA256 hash to brute force
  -t int
        Number of goroutines (default 4)
  -w string
        Newline delimited wordlist
```

### Outline

The program takes the hash target via command line flag `-h` in hex encoded form
The program reads words from a new-line delimited wordlist specified with `-w` 
Worker go routines are started, default of `4` but can be specified by `-t` 
Each word is sent to a channel read by workers that hashes and compare them to the target hash
Reading of the wordlist, and hashing attempts are stopped via context cancellation if a match is found.