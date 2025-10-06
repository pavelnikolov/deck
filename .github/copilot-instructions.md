# deck Development Guidelines

Last updated: 2025-10-06

## Active Technologies

- Go 1.25+ + None (stdlib only - crypto/rand, encoding/binary, sort).
- Avoid external dependencies.

## Project Structure

```
/deck.go          # Core types (Deck, Card, Suit, Rank) + joker support + Must* methods
/deck_test.go     # All tests including new joker tests and Must* panic tests
/example_test.go  # Example code including joker examples and Must* examples
/go.mod           # Module definition
```

## Commands

# Add commands for Go 1.25+
- `go test -v ./...` - Run all tests with verbose output
- `go test -cover ./...` - Check test coverage
- `go test -bench=. ./...` - Run all benchmarks
- `go test -run Example ./...` - Run all example tests
- `go fmt ./...` - Format all Go code
- `golangci-lint run ./...` - Run linter if golangci-lint is installed
- `go doc ./...` - View documentation for all packages

## Code Style

Go 1.25+: Follow standard conventions:
 - Effective Go https://go.dev/doc/effective_go
 - Go Proverbs https://go-proverbs.github.io
 - Go Code Review Comments https://go.dev/wiki/CodeReviewComments
 - Go Test Comments https://go.dev/wiki/TestComments
 - Google Go Style Guide https://google.github.io/styleguide/go/decisions
