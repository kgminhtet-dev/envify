# envify
.env Loader

## Getting started

### Getting envify

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/mr-kmh/envify"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Run the following command to install the envify package

```sh
$ go get -u github.com/mr-kmh/envify
```

### Running envify

```go
package main

import (
	"fmt"
	"os"

	"github.com/mr-kmh/envify"
)

func init() {
	envify.Load()
}

func main() {
	fmt.Println(os.Getenv("VARIABLE"))
}
```

Use Go command to run the demo:

```
$ go run example.go
```