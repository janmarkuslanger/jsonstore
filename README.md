# jsonstore

`jsonstore` is a tiny, file-backed key-value store for Go that persists data as human-readable JSON. It uses Go generics so you get strong typing for the values you read and write, and it stays concurrency-safe with an internal `sync.RWMutex`.

## Installation

```bash
go get github.com/janmarkuslanger/jsonstore
```

`jsonstore` targets Go 1.24 or newer (see `go.mod`).

## Quick Start

```go
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/janmarkuslanger/jsonstore"
)

type Profile struct {
	Name string
	Age  int
}

func main() {
	store, err := jsonstore.NewStore[Profile]("profiles.json")
	if err != nil {
		log.Fatalf("opening store failed: %v", err)
	}

	if err := store.Set("alice", Profile{Name: "Alice", Age: 31}); err != nil {
		log.Fatalf("store value failed: %v", err)
	}

	profile, err := store.Get("alice")
	if errors.Is(err, jsonstore.ErrNotFound) {
		log.Println("no profile stored yet")
		return
	} else if err != nil {
		log.Fatalf("read value failed: %v", err)
	}

	fmt.Printf("%s is %d years old\n", profile.Name, profile.Age)
}
```

## API Overview
- `jsonstore.NewStore[T](path string) (*Store[T], error)`: Instantiates a store bound to a JSON file. The file is created if it does not exist.
- `(*Store[T]).Set(key string, value T) error`: Serializes the value as JSON, stores it under the key, and persists the update.
- `(*Store[T]).Get(key string) (T, error)`: Returns the stored value or `jsonstore.ErrNotFound` if the key is missing.
- `(*Store[T]).Delete(key string) error`: Removes the key and persists the change.
- `(*Store[T]).Keys() ([]string, error)`: Lists all stored keys at the time of the call.

Internally, writes are serialized through a mutex and flushed to disk via a temporary file + rename dance to avoid partial files after crashes.

## Testing

Run the unit tests with:

```bash
go test ./...
```

## License

This project is distributed under the terms of the [MIT License](LICENSE).
