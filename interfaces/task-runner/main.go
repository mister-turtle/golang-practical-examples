package main

import (
	"fmt"
)

// this is the interface we're focusing on in this example
// interfaces should generally be narrow, defining as few functions as possible
type TaskRunner interface {
	Run() error
}

// we create arbitrary structs representing a task
// this disk struct takes a target and an operation
type DiskTask struct {
	Target    string
	Operation string
}

// by implementing this method we are implicitly satisfying the TaskRunner interface
func (d DiskTask) Run() error {
	fmt.Printf("Performing disk action %s on target %s\n", d.Operation, d.Target)
	return nil
}

// this http task struct contains fields relevant to an http operation
type HttpTask struct {
	URL            string
	Method         string
	ExpectedStatus int
}

// but we still implement Run() so that it can be used by something expecting a TaskRunner
func (h HttpTask) Run() error {
	return fmt.Errorf("Querying %s, expecting status %d, got %d\n", h.URL, h.ExpectedStatus, 302)
}

func main() {

	fmt.Println("Simple Interfaces Example")

	// create instances of our tasks with the expected parameters
	diskTask := DiskTask{Target: "/mnt/drive/somefolder", Operation: "DELETE"}
	httpTask := HttpTask{URL: "https://www.example.com", Method: "GET", ExpectedStatus: 200}

	// RunTasks has a variadic parameter, so we can call it with as many tasks as required
	// as long as they satisfy the TaskRunner interface
	err := RunTasks(diskTask, httpTask)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

}

// RunTasks takes a variadic parameter, which will be available to us as a slice of that type.
func RunTasks(tasks ...TaskRunner) error {

	// quick check to make sure we didn't receive zero arguments
	if len(tasks) == 0 {
		return fmt.Errorf("no tasks passed to the task runner")
	}

	// range over the slice of TaskRunners passed into the function
	for i, task := range tasks {

		fmt.Printf("\nRunning task %d\n", i)

		// the original types were passed in as a TaskRunner interface, which we know they satisfy at compile time
		err := task.Run()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

	}
	return nil
}
