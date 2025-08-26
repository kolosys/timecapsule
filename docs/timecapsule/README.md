# timecapsule



## Installation

```bash
go get github.com/kolosys/timecapsule/timecapsule
```

## Quick Start

```go
package main

import "github.com/kolosys/timecapsule/timecapsule"

func main() {
    // Your code here
}
```

## API Reference
### Functions
- [Example_basic](api-reference.md#example_basic) - Example_basic demonstrates basic usage of the time capsule

- [Example_context](api-reference.md#example_context) - Example_context demonstrates using context for cancellation

- [Example_management](api-reference.md#example_management) - Example_management demonstrates capsule management operations

- [Example_struct](api-reference.md#example_struct) - Example_struct demonstrates using structs with the time capsule

- [Example_waitForUnlock](api-reference.md#example_waitforunlock) - Example_waitForUnlock demonstrates waiting for a capsule to unlock

### Types
- [Capsule](api-reference.md#capsule) - Capsule represents a time-locked value

- [Codec](api-reference.md#codec) - Codec defines how to serialize/deserialize values

- [JSONCodec](api-reference.md#jsoncodec) - JSONCodec implements Codec using JSON encoding

- [MemoryTimeCapsule](api-reference.md#memorytimecapsule) - MemoryTimeCapsule implements TimeCapsule using in-memory storage

- [Metadata](api-reference.md#metadata) - Metadata contains information about a capsule without exposing its value

- [PersistentTimeCapsule](api-reference.md#persistenttimecapsule) - PersistentTimeCapsule implements TimeCapsule using a persistent storage backend

- [Storage](api-reference.md#storage) - Storage defines the interface for persistent storage backends

- [TimeCapsule](api-reference.md#timecapsule) - TimeCapsule is the main interface for storing and retrieving time-locked values


## Examples

See [examples](examples.md) for detailed usage examples.
