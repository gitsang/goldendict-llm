# Agent Development Guidelines for goldendict-llm

## Project Overview
GoldenDict external program written in Go for translating with LLM API. Single-package application using Cobra CLI and template-based rendering for dictionary lookup and sentence translation.

## Build, Lint, and Test Commands

### Build
```bash
go build -o goldendict-llm .
```

### Run Tests
```bash
# Run all tests in current package
go test -v .

# Run a specific test function
go test -v -run TestFunctionName .

# Run tests with coverage
go test -v -cover .
```

### Lint and Format
```bash
# Format code (run before committing)
go fmt .

# Run go vet for static analysis
go vet .

# Check formatting without writing
go fmt -n .
```

### Clean
```bash
# Clean build artifacts
go clean
```

## Code Style Guidelines

### Imports
- Group imports: standard library first, then third-party packages
- Use blank imports only for side effects (e.g., `_ "embed"`)
- Example order: stdlib → external → internal
- Keep imports alphabetical within groups

### Formatting
- Use `gofmt` standard formatting (tabs for indentation)
- No trailing whitespace
- Maximum line length: follow Go conventions (typically 100-120 chars)
- Use blank lines between functions and top-level declarations

### Types and Structs
- Use PascalCase for exported types and struct fields
- Use camelCase for unexported fields
- Use struct tags for JSON serialization and default values
- Example:
  ```go
  type AdapterConfig struct {
    Name  string
    URL   string
    Token string
    Model string
  }
  ```

### Naming Conventions
- **Functions/Methods**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase
- **Constants**: PascalCase for exported, camelCase for unexported
- **Interfaces**: Typically PascalCase, often ending in "er" suffix
- **File names**: lowercase_with_underscores or simple lowercase (e.g., `translator.go`, `template_word.go`)

### Error Handling
- Always check errors, never ignore them
- Wrap errors with context using `fmt.Errorf("operation failed: %v", err)`
- Use `panic` only in unrecoverable situations (as seen in main.go)
- Return errors from functions, don't swallow them
- Example:
  ```go
  if err != nil {
    return "", fmt.Errorf("json marshal failed: %v", err)
  }
  ```

### Functions and Methods
- Keep functions focused and short
- Use the functional options pattern for constructors (see `NewTranslator`)
- Document exported functions if behavior is non-obvious
- Use receiver names that are consistent (e.g., `t` for Translator)
- Example:
  ```go
  type TranslatorOption func(*Translator)
  
  func WithHTTPClient(client *http.Client) TranslatorOption {
    return func(t *Translator) {
      t.client = client
    }
  }
  ```

### Constants and Variables
- Define constants for string literals and configuration values
- Use `const` for compile-time constants
- Use `var` for variables with default struct tags
- Use type aliases for string constants (e.g., `type TranslateTemplate string`)

### Templates and Static Files
- Use `//go:embed` directive to embed static files
- Define template variables at package level for reuse
- Use `template.Must` for template parsing to panic on errors
- Register custom functions with `template.FuncMap`
- Example:
  ```go
  //go:embed static/word-prompt.md
  var WordPrompt string
  
  var wordTemplate = template.Must(template.New("word.html.tmpl").Parse(wordTemplateString))
  ```

### Testing
- Test files should end with `_test.go`
- Use table-driven tests for multiple test cases when appropriate
- Test functions should start with `Test`
- Use `t.Fatal` for failures that should stop test execution
- Use `t.Logf` for debugging information during tests
- Keep test helpers unexported if possible

### JSON and Configuration
- Use JSON struct tags for API models
- Use `omitempty` for optional fields
- Define separate structs for Request and Response models
- Example:
  ```go
  type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
  }
  ```

### HTTP Client Usage
- Accept `*http.Client` as a parameter or use functional options
- Set timeouts on HTTP clients
- Always close response bodies using defer
- Check HTTP status codes before processing response
- Example:
  ```go
  resp, err := t.client.Do(req)
  if err != nil {
    return "", fmt.Errorf("API request failed: %v", err)
  }
  defer func() {
    _ = resp.Body.Close()
  }()
  ```

### File Organization
- Single-package structure (all files in `main` package)
- Group related functionality in separate files (e.g., `translator.go`, `model.go`, `template_word.go`)
- Keep `main.go` for CLI setup and entry point
- Use separate files for different template types

### Comments
- Comments are minimal in this codebase
- Use package-level comments to explain file purpose
- Use inline comments only for complex logic
- No godoc comments for most exported functions

## Notes
- This is a simple CLI application with a single package
- No external linting configuration files (golangci-lint, etc.)
- Tests currently may be outdated and need fixing before running
- The application uses embed directives, so static files must exist at build time
