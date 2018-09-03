# Smartling SDK in Golang [![](https://godoc.org/github.com/Smartling/api-sdk-go?status.svg)](http://godoc.org/github.com/Smartling/api-sdk-go)

## Examples

Examples are located in the `example_*_test.go` files.

Examples suited for use with real user accounts, so to run examples you need
at least obtain your User ID and Token Secret.

Then, obtained User ID, Token Secret and other parameters should be populated
in the `example_credentials_test.go` file.

To run all examples, just run `go test`.

To run specific example it must be specified like `go test -run Projects_List`
(will run example from `examples_projects_list_test.go`).
