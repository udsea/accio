# Contributing to Accio

Thank you for your interest in contributing to Accio! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and considerate of others when contributing to this project.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with the following information:

- A clear, descriptive title
- Steps to reproduce the bug
- Expected behavior
- Actual behavior
- Any relevant logs or screenshots
- Your environment (OS, Go version, etc.)

### Suggesting Features

If you have an idea for a new feature, please create an issue with:

- A clear, descriptive title
- A detailed description of the feature
- Any relevant examples or mockups
- Why this feature would be useful

### Pull Requests

1. Fork the repository
2. Create a new branch for your changes
3. Make your changes
4. Write tests for your changes
5. Run the tests to make sure they pass
6. Submit a pull request

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/accio.git
   cd accio
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build ./cmd/accio
   ```

4. Run tests:
   ```bash
   go test ./...
   ```

## Adding New Sites

To add a new site to Accio:

1. Open `internal/sites/sites.go`
2. Add a new entry to the `GetSites()` function:
   ```go
   {
       Name:        "SiteName",
       URL:         "https://example.com/{}",
       ErrorType:   "status_code",
       URLProbe:    false,
       URLFormat:   "https://example.com/{}",
       CheckMethod: "GET",
   },
   ```
3. Test that the site works correctly
4. Submit a pull request

## Code Style

- Follow standard Go code style and conventions
- Use `gofmt` to format your code
- Write clear, descriptive comments
- Write meaningful commit messages

## Testing

- Write tests for new features and bug fixes
- Make sure all tests pass before submitting a pull request
- Run tests with `go test ./...`

## Documentation

- Update documentation for new features or changes
- Keep the README up to date
- Document public functions and types

## License

By contributing to Accio, you agree that your contributions will be licensed under the project's MIT license.