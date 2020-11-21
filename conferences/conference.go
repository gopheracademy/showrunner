package conferences

import "context"

type FooParams struct {
	Name string
}

type FooResponse struct {
	Message string
}

// Foo is an example endpoint.
// encore:api public
func Foo(ctx context.Context, params *FooParams) (*FooResponse, error) {
	message := "Hello, " + params.Name
	return &FooResponse{Message: message}, nil
}
