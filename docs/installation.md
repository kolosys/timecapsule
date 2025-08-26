# Installation

## Requirements

- Go 1.22 or later
- No external dependencies required

## Install via go get

```bash
go get github.com/kolosys/timecapsule@latest
```

## Install specific version

```bash
go get github.com/kolosys/timecapsule@v0.1.0
```

## Verify installation

Create a simple test file:

```go
package main

import (
    "fmt"
    
    "github.com/kolosys/timecapsule"
)

func main() {
    fmt.Println("timecapsule installed successfully!")
}
```

Run it:

```bash
go run main.go
```

## Module integration

Add to your `go.mod`:

```bash
go mod init your-project
go get github.com/kolosys/timecapsule@latest
```

## Import packages

```go
import "github.com/kolosys/timecapsule"
```
