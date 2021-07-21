# Golang Practical Examples

## Overview

This is a basic repository where I put practical examples that demonstate core language concepts that follow idiomatic usage of Go.

## Concurrency Examples
[cross-routine](/concurrency/cross-routine) - Number counter that shares channels across goroutines  
[portscan-pipeline](/concurrency/tcp-portscan) - TCP port scanner using fan-out pattern  
[sha256-bruteforcer](/concurrency/sha256-bruteforcer) - SHA256 brute forcer using fan-out pattern, and context cancellation  

## HTTP Examples
[simplest-server](/http/simplest-server) - HTTP server using stdlib net/http package with custom headers, cookies, and redirects  
[simplest-templates](/http/simplest-templates) - Easy introduction to templates using the stdlib  

## Interface Examples
[task-runner](/interfaces/task-runner) - Execute a variety of arbitrary tasks via a common task executor using interfaces  
