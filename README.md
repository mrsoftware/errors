# errors
> [!WARNING]  
> Be aware this package is not ready for production usage and is in the development phase.

Package errors help with adding context to your errors.

```shell
go get github.com/mrsoftware/errors
```

The traditional error handling idiom in Go is roughly akin to
```go
if err != nil {
    return err
}
```

which applied recursively up the call stack results in error reports without context or debugging information. 
The errors package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error.

The errors package supports various context adding capabilities like other error packages,
but the important part and the thing make it unique, is the ability to add multiple data in any type.
some time you need to add contextual data into your error, in this case information about the user
```go
user := struct { Username string }{Username: "mrsoftware"}
errors.New("some error", errors.String("name", "mohammad"), errors.Reflect("user", user))
```

to retrieve field from a single error
```go
err := errors.New("some error", errors.String("username", "mrsoftware"))
nameField := errors.GetField(err, "username")
```

but if you need to look for field in the error chain
```go
cause := errors.New("cause error", errors.String("username", "mrsoftware"))
err := errors.Wrap(err, "some error", errors.String("name", "mohammad"))
nameField := errors.FindFieldInChain(err, "username")
```

or you can wrap the cause error:

```go
_, err := ioutil.ReadAll(r)
if err != nil {
        errors.Wrap(err, "some error", errors.String("name", "mohammad"))
}
```

as errors support the `causer` interface you can get the cause:

```go
type causer interface {
        Cause() error
}
```

`errors.Cause` recursively try to find the error that supports the causer and retrieve the cause.

for mode details, check the [documentation](https://godoc.org/github.com/mrsoftware/errors)