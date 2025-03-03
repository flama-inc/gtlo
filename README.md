# Go Tiny Lock Object

## Build

- `task build`

## Usage

```go
lock := tlo.New("/path/to")

if err := lock.Load(); err != nil {
    // error
}

if !lock.IsLocked() {
    lock.Lock()
    lock.Save()
    // e.g. return error
}
```

## Metadata

```go
lock.SetMetadata("key", []byte("value"))

b, err := lock.GetMetadata("key")
// []byte("value")

bs, err := lock.GetMetadataAll()
// map[string][]byte
```
