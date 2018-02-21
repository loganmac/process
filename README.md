# Process

A simple external process manager for go. Create a processor by creating a `Driver`, and calling `process.New(driver)` with your driver.

Then call `processor.Run(taskName, cmd, args...)` to run external commands with the processor.

Will execute the process, initializing your driver with the task name, and respond to stdoud, stderr, and exit from the process.

Works on Unix and Windows (though the example in the folder is specific to Unix)

## Example

```go
import "github.com/loganmac/process"

func main() {
  driver := &Summarizer{} // something that implements process.Driver
  processor := process.New(driver)
  if err := processor.Run("list all files", "ls", "l", "a", "h"); err != nil{
    log.Fatal(err)
  }
}
```

For a more complex example, see [the example folder](example).

It shows how to use this library (and some others) to build this:

![example-gif](./example.gif)
