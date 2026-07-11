# Contributing to OmniConfig

Thank you for considering contributing! We welcome contributions of all kinds.

## Code of Conduct

Be respectful, inclusive, and constructive.

## How to Contribute

1. **Fork** the repository on GitHub
2. **Create a feature branch** — `git checkout -b feat/your-feature`
3. **Write your code** with tests
4. **Run tests** — `make test`
5. **Format your code** — `make fmt`
6. **Run linters** — `make lint`
7. **Commit** — use clear, descriptive commit messages
8. **Push** — `git push origin feat/your-feature`
9. **Open a pull request**

## Development Setup

```bash
# Install Go 1.22+
git clone https://github.com/YOUR_USERNAME/cli.git
cd cli
make build
make test
```

## Pull Request Guidelines

- Keep PRs focused — one feature/fix per PR
- Write descriptive titles and summaries
- Include tests for new functionality
- Update documentation if needed
- Ensure all tests pass

## Adding a New Config Format

1. Create a new file under `pkg/formats/`
2. Implement the `Handler` interface
3. Register in `init()`
4. Add tests

## Mobile Platform Support

- **iOS**: Build with `GOOS=ios GOARCH=arm64`. Requires Xcode command line tools.
- **Android**: Build with `GOOS=android GOARCH=arm64|amd64`. Can use NDK for cgo builds.

## License

By contributing, you agree your contributions will be licensed under the MIT License.