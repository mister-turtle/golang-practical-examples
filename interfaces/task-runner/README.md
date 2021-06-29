# Task Runner
  
### Overview
This example shows using an interface to describe the functionality of a struct.  
A function that accepts an interface will accept any struct that implements that interface.  

In this example we have tasks that perform different actions, we have:
 * Disk Task
 * HTTP Task

We would like these to be run by a Task Runner, a function that can orchestrate the running of different tasks.  
To keep it simple, all the task runner does is execute the task, however in a real scenario this could involve concurrency, some distributed system, mechanisms for retrying or re-scheduling etc.  
Also to keep it simple, the tasks themselves just print a message.

An HTTP task requires different parameters, and will perform a very different function than the Disk task. We could abstract some type of struct that encompasses the type of and number of arguments for both, but this gets convoluted and limiting very quickly. Instead, we define a TaskRunner interface, which has a single method `Run()`.

This allows us to create a function that takes in any struct that satisfies the `TaskRunner` interface, and execute it's `Run()` method, regardless of what the underlying type is. Interfaces describe functionality, not properites. So we use different struct types to hold the relevant parameters for each task, and make sure they satisfy the `TaskRunner` interface.

### Running

```shell
go run main.go

Simple Interfaces Example

Running task 0
Performing disk action DELETE on target /mnt/drive/somefolder

Running task 1
Performing http action GET on https://www.example.com expecting 200
Error: Querying https://www.example.com, expecting status 200, got 302
```
