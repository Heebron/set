# set

A small, generic Set implementation for Go that supports both standard and concurrency-safe variants.

- Generic over any comparable type (Go 1.18+ generics)
- Optional concurrency-safe implementation for use across multiple goroutines
- Simple, focused API for common set operations
- Zero external dependencies

## Module

Module path: `github.com/Heebron/set/v2`

Go version: 1.21+

## Installation

```bash
go get github.com/Heebron/set/v2
```

## Quick Start

```go
package main

import (
    "fmt"
    set "github.com/Heebron/set/v2"
)

func main() {
    // Non-concurrent set (faster, not safe for concurrent access)
    s := set.New[string]()
    s.Add("apple")
    s.Add("banana")

    fmt.Println(s.Contains("apple")) // true
    fmt.Println(s.Size())             // 2

    // Concurrent set (safe for access from multiple goroutines)
    cs := set.NewConcurrent[int]()
    cs.Add(1)
    cs.Add(2)

    // Wait for the set to become empty (or time out)
    emptied := cs.WaitForEmptyWithTimeout(0) // immediate check
    fmt.Println(emptied) // false, because it isn't empty yet

    cs.Remove(1)
    cs.Remove(2)
    emptied = cs.WaitForEmptyWithTimeout(0)
    fmt.Println(emptied) // true
}
```

## API Overview

The package exposes a single generic type `Set[T comparable]` and a small set of constructors and methods.

Constructors:
- `New[T comparable]() Set[T]`
  - Creates a non-concurrent set.
- `NewWithInitializer[T comparable](members ...T) Set[T]`
  - Creates a non-concurrent set pre-populated with the provided members.
- `NewConcurrent[T comparable]() Set[T]`
  - Creates a concurrency-safe set guarded by an `RWMutex`.
- `NewConcurrentWithInitializer[T comparable](members ...T) Set[T]`
  - Creates a concurrency-safe set pre-populated with the provided members.

Core methods:
- `(*Set[T]) Add(e T) bool`
  - Inserts an element. Returns true if the set changed (element was new).
- `(*Set[T]) Remove(e T) bool`
  - Removes an element. Returns true if the set changed (element existed).
- `(*Set[T]) Contains(e T) bool`
  - Reports whether the element is present.
- `(*Set[T]) Size() int`
  - Returns the number of elements.
- `(*Set[T]) IsEmpty() bool`
  - Convenience: `Size() == 0`.
- `(*Set[T]) Members() []T`
  - Returns a snapshot slice of all members; order is unspecified.
- `(*Set[T]) Clear()`
  - Removes all elements; preserves the concurrency mode.
- `(*Set[T]) Clone() Set[T]`
  - Returns a shallow copy in the same concurrency mode.

Set algebra:
- `(*Set[T]) Union(rhs Set[T]) Set[T]`
- `(*Set[T]) Intersect(rhs Set[T]) Set[T]`
- `(*Set[T]) Difference(rhs Set[T]) Set[T]`
- `(*Set[T]) IsSubset(rhs Set[T]) bool`
- `(*Set[T]) Equal(rhs Set[T]) bool`

Concurrency utility:
- `(*Set[T]) WaitForEmptyWithTimeout(timeout time.Duration) bool`
  - Only valid on concurrent sets; panics if called on a non-concurrent set.
  - Waits until the set transitions to empty or the timeout elapses, then returns whether the set is empty at that moment.

## Usage Examples

Creating and combining sets:

```go
s1 := set.NewWithInitializer(1, 2, 3)
s2 := set.NewWithInitializer(3, 4)

u := s1.Union(s2)      // {1,2,3,4}
i := s1.Intersect(s2)  // {3}
d := s1.Difference(s2) // {1,2}

fmt.Println(u.Size(), i.Size(), d.Size())
```

Concurrent usage:

```go
cs := set.NewConcurrent[string]()
var wg sync.WaitGroup

wg.Add(2)

go func() {
    defer wg.Done()
    cs.Add("a")
    cs.Add("b")
}()

go func() {
    defer wg.Done()
    _ = cs.Contains("a")
    cs.Remove("a")
}()

wg.Wait()
_ = cs.WaitForEmptyWithTimeout(time.Millisecond) // check if empty soon
```

## Design Notes

- The set is backed by a `map[T]struct{}` equivalent (an empty value type) to minimize memory overhead.
- Concurrency-safe sets embed an `*sync.RWMutex` to guard all operations. Non-concurrent sets omit the lock to avoid overhead.
- Methods that return new sets (Union, Intersect, Difference, Clone) preserve the concurrency mode of the receiver.
- The zero value of `Set[T]` is not usable. Always use one of the constructors.

## Performance Considerations

- Prefer `New[T]` for single-goroutine scenarios to avoid locking overhead.
- Use `Members()` to obtain a snapshot; the order is intentionally unspecified and may vary between calls.
- All read methods on concurrent sets use read locks to allow parallelism where possible.

## Errors and Panics

- Calling `WaitForEmptyWithTimeout` on a non-concurrent set will panic. Use `NewConcurrent`/`NewConcurrentWithInitializer` if you need this behavior.
- Timeouts are treated as non-negative; negative durations are coerced to zero (immediate check).

## Versioning and Compatibility

- Requires Go 1.21 or later (as per `go.mod`).
- Uses generics and thus requires Go 1.18+; tested and intended for 1.21+.

## Testing

Run the tests with:

```bash
go test ./...
```

## License

This project is licensed under the terms of the MIT License. See [LICENSE](LICENSE).

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) and our [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
