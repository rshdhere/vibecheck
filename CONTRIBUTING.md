# Contributing to Vibecheck

Thank you for your interest in contributing to Vibecheck! This document provides guidelines and instructions for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/vibecheck.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit your changes: `git commit -m "Add your meaningful commit message"`
7. Push to your fork: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git
- Access to at least one LLM provider API key for testing

### Building the Project

```bash
# Clone the repository
git clone https://github.com/rshdhere/vibecheck.git
cd vibecheck

# Build the project
go build -o vibecheck

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Running Locally

```bash
# Install dependencies
go mod download

# Run the application
./vibecheck commit

# Or build and install
go install
```

## Code Style

- Follow Go conventions and best practices
- Use `gofmt` to format your code
- Follow the existing code style in the project
- Write clear, self-documenting code
- Add comments for complex logic

### Formatting

```bash
# Format your code
go fmt ./...

# Run go vet
go vet ./...
```

## Testing

- Write tests for new features
- Ensure all existing tests pass
- Aim for good test coverage
- Test with multiple LLM providers when possible

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Commit Messages

- Use clear, descriptive commit messages
- Follow conventional commit format when possible
- Reference issue numbers if applicable
- Keep commits focused on a single change

Example:
```
feat: add support for new LLM provider
fix: resolve API key validation issue
docs: update installation instructions
```

## Pull Request Process

1. Ensure your code follows the project's style guidelines
2. Update documentation if needed
3. Add tests for new features
4. Ensure all tests pass
5. Update CHANGELOG.md if applicable
6. Write a clear PR description explaining:
   - What changes were made
   - Why the changes were made
   - How to test the changes
   - Any breaking changes

### PR Checklist

- [ ] Code follows the project's style guidelines
- [ ] Tests have been added/updated
- [ ] All tests pass
- [ ] Documentation has been updated
- [ ] Commit messages are clear and descriptive
- [ ] No merge conflicts with main branch

## Adding New LLM Providers

If you want to add support for a new LLM provider:

1. Create a new client file in `internal/llm/provider-name/client.go`
2. Implement the `Provider` interface from `internal/llm/provider.go`
3. Add the provider to the provider registry
4. Update documentation
5. Add tests for the new provider
6. Update the README with the new provider

## Reporting Bugs

When reporting bugs, please include:

- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)
- Relevant error messages or logs
- Any additional context

## Suggesting Features

When suggesting features:

- Provide a clear description of the feature
- Explain the use case and benefits
- Consider implementation complexity
- Check if a similar feature already exists
- Be open to discussion and feedback

## Code Review

- Be respectful and constructive in reviews
- Focus on the code, not the person
- Provide specific, actionable feedback
- Acknowledge good work
- Be open to feedback on your own code

## Questions?

If you have questions about contributing:

- Open an issue with the "question" label
- Check existing issues and discussions
- Review the documentation

## License

By contributing to Vibecheck, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in the project's documentation and release notes. Thank you for helping make Vibecheck better!

