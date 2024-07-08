# errors
## Enriching Error Handling in Go

The errors package aims to empower you with a more robust approach to error handling in your Go applications. It addresses the limitations of Go's default error handling by providing features for:

- **Contextual Information**: Add meaningful details to errors, making them more informative for debugging and troubleshooting.
- **Error Chaining**: Create a chain of errors, capturing the root cause and subsequent events leading to the final error.
- **Multiple Data**: Attach various types of data (strings, structs, etc.) to errors for richer context.
- **MultiError**: Group and manage multiple errors without losing their individual structure. It's concurrent safe, can be used in multiple goroutine.
- **WaitGroup**: Streamline m error handling by potentially providing abstractions for common use cases like MultiError + sync.WaitGroup.

## Installation

```shell
go get github.com/mrsoftware/errors
```

## Beyond Basic Error Handling
Go's traditional error handling often results in error messages lacking context and debugging information, especially when errors propagate up the call stack. The errors package tackles this challenge by allowing you to enrich error messages with valuable data.

## Adding Context
The errors.New function is your go-to tool for creating errors with context. You can pass additional arguments to specify different types of data:
```go
user := struct { Username string }{Username: "mrsoftware"}
errors.New("some error", errors.String("name", "mohammad"), errors.Reflect("user", user))
```
 
## Retrieving Context
To extract the added data from an error:

- Single Error:
```go
err := errors.New("some error", errors.String("username", "mrsoftware"))
nameField := errors.GetField(err, "username")
```

- Error Chain:
```go
cause := errors.New("cause error", errors.String("username", "mrsoftware"))
err := errors.Wrap(err, "some error", errors.String("name", "mohammad"))
nameField := errors.FindFieldInChain(err, "username")
```

## Wrapping Existing Errors
The errors.Wrap function lets you add context to existing errors while preserving the original cause:
```go
_, err := ioutil.ReadAll(r)
if err != nil {
        errors.Wrap(err, "some error", errors.String("name", "mohammad"))
}
```

## Retrieving Cause
The package aims to incorporate a `Cause()` function that will assist you in traversing the error chain and retrieving the root cause error.


## Handling Multiple Errors
In scenarios where you encounter multiple errors concurrently, the `errors.MultiError` type offers a seamless way to aggregate them without resorting to cumbersome string concatenation. It also ensures thread safety for concurrent operations:

```go

err := errors.NewMultiError(io.EOF, errors.New("getting data"))

// Real Usage: when calling a third party service I need to check if error is timeout signal the up layer and do not want them notify timeout using lower level code, but my timeout. 
err := callingHttpCode()
if os.IsTimeout(errors.Cause(err)) {
    return res, errors.NewMultiError(err, constants.ErrTimeout)
}
```
**Real Usage**: some time happened you need to create list of error from multiple goroutine.
```go
//
errorList := errors.NewMultiError()
wg := &sync.WaitGroup{}

wg.Add(1)
go func(){
	err := callingHttpClient()
	errorList.SafeAdd(err) // if result of `callingHttpClient` is nil no error will add to the err list.

	wg.Done()
}()

wg.Add(1)
go func(){
    err := callingHttpClient()
	if err != nil {
        errorList.SafeAdd(err)
    }

	wg.Done()
}()

wg.Wait()

if err := errorList.Err(); err != nil {
	// oh, something bad happened in one of routines above.
}
```

## Combining sync.WaitGroup and errors.MultiError
The `errors.WaitGroup` type simplifies error handling in concurrent operations by merging `sync.WaitGroup` with `MultiError`. This reduces boilerplate code for cleaner and more concise error handling:
```go
wg := errors.NewWaitGroup() // you can pass some options, e.g: your custom sync.WaitGroup using WaitGroupWithSyncWaitGroup.

wg.Add(1)
go func(){
    wg.Done(callingHttpClient())
}()

wg.Add(1)
go func(){
    err := callingHttpClient()
	wg.Done(err)
}()



if err := wg.Wait(); err != nil {
    // oh, something bad happened in one of routines above.
} 
```

**or you can use Do method and let WaitGroup handle Add and Done internally.**

```go
wg := errors.NewWaitGroup() 

wg.Do(func() error {
    return callingHttpClient()
})

wg.Do(func() error {
    return callingHttpClient()
})


if err := wg.Wait(); err != nil {
    // oh, something bad happened in one of routines above.
} 
```

## limiting the concurrent task counts using limit Options in errors.WaitGroup
```go
wg := errors.NewWaitGroup(errors.WaitGroupWithTaskLimit(2)) 

wg.Do(func() error {
    return callingHttpClient()
})

wg.Do(func() error {
    return callingHttpClient()
})

// we set limit concurrent task to 2, so this task will block until one of above are done.
wg.Do(func() error {
	return callingHttpClient()
})

if err := wg.Wait(); err != nil {
    // oh, something bad happened in one of routines above.
} 
```

## Use Custom runner instead of GoRoutine in errors.WaitGroup
```go
import (
    "github.com/mrsoftware/errors"
    "github.com/panjf2000/ants/v2"
)

// in this example we are using ants goroutine pool.
wg := errors.NewWaitGroup(errors.WaitGroupWithTaskRunner(ants.Submit)) 

wg.Do(func() error {
    return callingHttpClient()
})

wg.Do(func() error {
    return callingHttpClient()
})
 
if err := wg.Wait(); err != nil {
    // oh, something bad happened in one of routines above.
} 
```


for mode details, check the [documentation](https://godoc.org/github.com/mrsoftware/errors)

----
## Roadmap
- [ ] Unit test
- [ ] complete toolbox
- [X] Multi error
- [X] Waiting error (sync.Waiting + errors)