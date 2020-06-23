# Memory leaks detection test

## Case

Add a testing scenario to the `runner` function.

## Start

```bash
go run ./cmd/pprof -strategy saml
```

and wait some time, depending on a testing scenario.

If the memory usage is stable that's great.

## Profiling

Fetch memory heap:

```bash
curl -sK -v http://localhost:6868/debug/pprof/heap > heap.out
```

Check a profile snapshot:

```bash
go tool pprof heap.out
```

When type `web` and observe the memory heap graph for anything weird.