# Contributing

Thank you for your interest in contributing to Glance! This guide will help you get started.

## Code of Conduct

We are committed to providing a welcoming and inclusive environment. By participating, you agree to:

- Be respectful and considerate
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards other community members

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report:

1. **Check existing issues**: Search [GitHub Issues](https://github.com/wahyudotdev/glance/issues) to see if it's already reported
2. **Try latest version**: Verify the bug exists in the latest release
3. **Gather information**: Collect logs, screenshots, and steps to reproduce

Create a bug report including:

- **Title**: Clear, descriptive summary
- **Environment**: OS, Glance version, Go version
- **Steps to Reproduce**: Detailed steps to trigger the bug
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happens
- **Logs**: Relevant log output (use `--log-level debug`)
- **Screenshots**: If applicable

**Example**:

```markdown
**Title**: Traffic not appearing in dashboard for HTTPS requests

**Environment**:
- OS: macOS 14.0
- Glance: v0.1.4
- Go: 1.24.4

**Steps to Reproduce**:
1. Start Glance with `glance`
2. Set proxy: `export HTTPS_PROXY=http://localhost:15500`
3. Make request: `curl https://api.github.com/users`

**Expected**: Request appears in dashboard
**Actual**: Only HTTP requests visible, HTTPS missing

**Logs**:
```
2026-02-22 10:30:00 DEBUG Intercepted request...
```
```

### Suggesting Features

Feature requests are welcome! Before creating one:

1. **Check roadmap**: See if it's already planned
2. **Search issues**: Look for similar suggestions
3. **Consider scope**: Ensure it fits Glance's goals

Create a feature request including:

- **Use Case**: Why is this needed?
- **Proposed Solution**: How should it work?
- **Alternatives**: Other approaches considered
- **Examples**: How other tools handle this

### Pull Requests

#### Setup Development Environment

1. **Fork the repository**
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/glance.git
   cd glance
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/wahyudotdev/glance.git
   ```
4. **Install dependencies**:
   ```bash
   # Backend
   go mod download

   # Frontend
   cd web/dashboard && npm install
   ```

#### Making Changes

1. **Create a branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Follow the [Code Style Guide](#code-style)
   - Add tests for new features
   - Update documentation if needed

3. **Run tests**:
   ```bash
   # Backend tests
   make test

   # Frontend tests
   cd web/dashboard && npm test

   # Linting
   make lint
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add amazing feature"
   ```

   Follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `test:` - Adding tests
   - `refactor:` - Code refactoring
   - `perf:` - Performance improvements
   - `chore:` - Maintenance tasks

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**:
   - Go to [GitHub](https://github.com/wahyudotdev/glance/pulls)
   - Click "New Pull Request"
   - Select your branch
   - Fill out the PR template

#### Pull Request Checklist

Before submitting, ensure:

- [ ] Tests pass locally (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Code follows style guide
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] No breaking changes (or clearly documented)
- [ ] Screenshots for UI changes
- [ ] Changelog entry added (for significant changes)

#### Pull Request Review Process

1. **Automated Checks**: GitHub Actions runs tests
2. **Code Review**: Maintainers review your code
3. **Feedback**: Address any requested changes
4. **Approval**: Once approved, PR will be merged
5. **Release**: Changes included in next release

## Code Style

### Go Code

Follow [Effective Go](https://go.dev/doc/effective_go) guidelines:

**Good**:
```go
// HandleTraffic processes intercepted HTTP traffic
func HandleTraffic(req *http.Request) error {
    if req == nil {
        return errors.New("request cannot be nil")
    }

    // Clear, readable logic
    traffic := &Traffic{
        ID:     generateID(),
        Method: req.Method,
        URL:    req.URL.String(),
    }

    return storage.Save(traffic)
}
```

**Bad**:
```go
func handle(r *http.Request) error {  // Unclear name
    t := &Traffic{r.Method, r.URL.String()}  // Missing fields
    return storage.Save(t)  // No error checking
}
```

**Naming**:
- Use `camelCase` for private, `PascalCase` for public
- Be descriptive: `GetTrafficByID` not `GetT`
- Avoid stuttering: `traffic.Traffic` → `traffic.Entry`

**Error Handling**:
```go
// Return errors, don't panic
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### TypeScript/React

Follow [Airbnb React Style Guide](https://github.com/airbnb/javascript/tree/master/react):

**Good**:
```typescript
interface TrafficProps {
  traffic: Traffic[];
  onSelect: (id: string) => void;
}

export const TrafficList: React.FC<TrafficProps> = ({
  traffic,
  onSelect
}) => {
  return (
    <div className="traffic-list">
      {traffic.map(item => (
        <TrafficItem
          key={item.id}
          traffic={item}
          onClick={() => onSelect(item.id)}
        />
      ))}
    </div>
  );
};
```

**Component Structure**:
1. Imports
2. Type definitions
3. Component definition
4. Styles (if applicable)
5. Export

**Hooks**:
- Use custom hooks for reusable logic
- Prefix with `use`: `useTraffic`, `useRules`
- Keep hooks focused and composable

### Documentation

**Code Comments**:
```go
// Good: Explains WHY, not WHAT
// Use buffered channel to prevent blocking when no listeners
traffic := make(chan *Traffic, 100)

// Bad: Explains obvious
// Create a channel
traffic := make(chan *Traffic)
```

**Package Documentation**:
```go
// Package proxy implements a MITM proxy for intercepting HTTP/HTTPS traffic.
//
// The proxy uses dynamic certificate generation to decrypt HTTPS traffic
// and applies rules for mocking and breakpoint functionality.
package proxy
```

**API Documentation**:
- Document all exported functions, types, and constants
- Include examples for complex APIs
- Keep documentation up-to-date with code changes

## Testing

### Backend Tests

**Unit Tests**:
```go
func TestHandleTraffic(t *testing.T) {
    tests := []struct {
        name    string
        request *http.Request
        want    error
    }{
        {
            name:    "valid request",
            request: httptest.NewRequest("GET", "/", nil),
            want:    nil,
        },
        {
            name:    "nil request",
            request: nil,
            want:    errors.New("request cannot be nil"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := HandleTraffic(tt.request)
            if !errors.Is(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Integration Tests**:
```go
func TestProxyIntegration(t *testing.T) {
    // Setup
    proxy := NewProxy()
    defer proxy.Close()

    // Test
    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxy.URL),
        },
    }

    resp, err := client.Get("https://example.com")
    if err != nil {
        t.Fatal(err)
    }

    // Verify
    if resp.StatusCode != 200 {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }
}
```

### Frontend Tests

**Component Tests**:
```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { TrafficList } from './TrafficList';

test('renders traffic items', () => {
  const traffic = [
    { id: '1', method: 'GET', url: 'https://example.com' }
  ];

  render(<TrafficList traffic={traffic} onSelect={() => {}} />);

  expect(screen.getByText('GET')).toBeInTheDocument();
  expect(screen.getByText('https://example.com')).toBeInTheDocument();
});

test('calls onSelect when item clicked', () => {
  const onSelect = jest.fn();
  const traffic = [
    { id: '1', method: 'GET', url: 'https://example.com' }
  ];

  render(<TrafficList traffic={traffic} onSelect={onSelect} />);

  fireEvent.click(screen.getByText('GET'));
  expect(onSelect).toHaveBeenCalledWith('1');
});
```

## Documentation

### Updating Docs

Documentation is in the `docs/` folder using Markdown + Docsify.

**Adding a New Page**:

1. Create `docs/new-page.md`
2. Add to `docs/_sidebar.md`:
   ```markdown
   - [New Page](new-page.md)
   ```
3. Test locally:
   ```bash
   cd docs
   npx serve
   ```
4. Preview at `http://localhost:3000`

**Documentation Guidelines**:
- Use clear, simple language
- Include code examples
- Add screenshots for UI features
- Link to related pages
- Keep navigation logical

## Community

### Communication Channels

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: Questions, ideas, general discussion
- **Pull Requests**: Code contributions

### Getting Help

- 📖 Read the [documentation](/)
- 💬 [Ask a question](https://github.com/wahyudotdev/glance/discussions)
- 🐛 [Report a bug](https://github.com/wahyudotdev/glance/issues)

## Recognition

Contributors are recognized in:
- **README**: Contributors section
- **Changelog**: For significant contributions
- **Release Notes**: Major features and fixes

## License

By contributing to Glance, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to:
- Open a [discussion](https://github.com/wahyudotdev/glance/discussions)
- Ask in your pull request
- Reach out to maintainers

Thank you for contributing to Glance! 🎉
