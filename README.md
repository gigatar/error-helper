# errorhelper

`errorhelper` is a small, extensible Go package for standardised HTTP error response encoding. It provides a consistent JSON error format, sensible HTTP status code mappings, and full support for user-defined behaviour.

## Features

* Default mapping from common application errors to HTTP status codes
* Custom error-to-status mapping
* Fully pluggable encoder via an interface
* JSON error response output with `code` and `string` fields

## Installation

```bash
go get github.com/gigatar/errorhelper
```

## Usage

### Basic Example

```go
import (
	"context"
	"net/http"

	"github.com/gigatar/errorhelper"
)

func handler(w http.ResponseWriter, r *http.Request) {
	err := someFunction()
	if err != nil {
		encoder := errorhelper.NewDefaultEncoder(nil)
		encoder.Encode(r.Context(), errorhelper.ErrBadRequest, w)
		return
	}
}
```

### JSON Output

A response might look like this:

```json
{
  "error": {
    "code": 400,
    "string": "Bad Request"
  }
}
```

## Custom Error Mapping

You can supply your own mapping function to customise which HTTP status code is used:

```go
customMapper := func(err error) int {
	switch err {
	case myapp.ErrDatabase:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

encoder := errorhelper.NewDefaultEncoder(customMapper)
encoder.Encode(ctx, myapp.ErrDatabase, w)
```

## Full Custom Encoder

You can implement the `Encoder` interface to replace encoding logic entirely:

```go
type MyEncoder struct{}

func (m *MyEncoder) Encode(ctx context.Context, err error, w http.ResponseWriter) {
	http.Error(w, "Something went wrong", http.StatusTeapot)
}
```

## Provided Errors

The package includes common error constants:

* `ErrNotFound`
* `ErrBadRequest`
* `ErrUnauthorized`
* `ErrForbidden`
* `ErrDuplicate`
* `ErrNotImplemented`
* `ErrNoContent`
* `ErrTooLarge`

## License

GPLv3