# Contributing to deck

Thank you for your interest in contributing to the deck package! We welcome contributions from the community.

## Guidelines

Please follow these guidelines when contributing:

1. **Follow Go best practices and Effective Go guidelines**
   - Write idiomatic Go code
   - Follow the conventions outlined in [Effective Go](https://go.dev/doc/effective_go)
   - Use `gofmt` to format your code

2. **Add tests for new features**
   - Maintain test coverage above 95%
   - Write both unit tests and example tests
   - Include edge cases and error scenarios

3. **Update documentation and examples**
   - Add godoc comments for all exported types and functions
   - Include example tests that demonstrate usage
   - Update README.md if adding major features

4. **Code quality checks**
   - Run `go fmt` to format your code
   - Run `go vet` to check for common mistakes
   - Run `golangci-lint run` if available
   - Ensure all tests pass with `go test -v`

## Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Testing

Before submitting a pull request, please ensure:

```bash
# Run all tests
go test -v

# Check test coverage
go test -cover

# Run benchmarks
go test -bench=.

# Run examples
go test -run Example
```

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to foster an open and welcoming environment.

## Questions?

If you have questions or need help, feel free to open an issue for discussion.
