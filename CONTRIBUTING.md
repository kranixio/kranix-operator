# Contributing to kranix-operator

Thank you for your interest in contributing to kranix-operator!

## Development Setup

```bash
# Clone the repository
git clone https://github.com/kranix-io/kranix-operator
cd kranix-operator

# Install dependencies
go mod download

# Install controller-gen
go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

# Generate CRD manifests and deepcopy methods
controller-gen crd:trivialVersions=true rbac:roleName=kranix-operator-role paths="./..." output:crd:artifacts:config=config/crd/bases

# Run locally (requires a cluster)
KRANE_CORE_ADDRESS=localhost:50051 go run ./cmd/operator --kubeconfig ~/.kube/config
```

## CRD Development

When modifying CRD types in `api/v1alpha1/`:

1. Update the Go types with `// +kubebuilder` markers
2. Regenerate deepcopy methods and CRD manifests:
   ```bash
   controller-gen crd:trivialVersions=true rbac:roleName=kranix-operator-role paths="./..." output:crd:artifacts:config=config/crd/bases
   ```
3. Install the updated CRDs to your cluster:
   ```bash
   kubectl apply -f config/crd/bases/
   ```

## Testing

### Unit Tests

```bash
go test ./internal/...
```

### E2E Tests with envtest

```bash
go test ./tests/e2e/... -tags e2e
```

## Adding a New Controller

1. Create the CRD type in `api/v1alpha1/yourresource_types.go`
2. Generate manifests with `controller-gen`
3. Create the controller in `internal/controllers/yourresource_controller.go`
4. Add reconciliation logic in `internal/reconciler/reconciler.go`
5. Register the controller in `cmd/operator/main.go`
6. Add RBAC markers to the controller
7. Write unit tests for the controller
8. Write E2E tests using envtest

## Code Style

- Follow standard Go formatting: `go fmt ./...`
- Run linter: `golangci-lint run`
- Use meaningful variable names
- Add comments for exported functions

## Commit Messages

Follow conventional commits format:

- `feat: add new controller for X`
- `fix: resolve reconciliation loop issue`
- `docs: update CRD documentation`
- `test: add E2E tests for Y`

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Commit your changes
6. Push to the branch
7. Open a Pull Request

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
