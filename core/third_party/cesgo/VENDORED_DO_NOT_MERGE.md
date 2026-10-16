# Temporary copy of the CES Go SDK: do not merge

> [!WARNING]
> TODO: Delete this directory and the `replace` directive in `core/go.mod`
> once the CES Go SDK is published, and depend on the published module
> instead. This copy exists only for early development and must not be merged.

This is a copy of the generated CoreWeave Evaluations Service (CES) Go SDK from
`services/evaluations/generated/go` in wandb/core at commit
`8726e0a110368c7871151b8ab3ed59e684c942ac`, without its `examples/`, `scripts/`,
and `Brewfile`.

wandb/core is private and wandb/wandb is public, so pushing a branch that
contains this directory publishes the SDK source.

To refresh it from a wandb/core checkout at `$CORE`:

```bash
git -C "$CORE" archive <commit> services/evaluations/generated/go | tar -x -C /tmp/cesgo
rsync -a --delete --exclude examples --exclude scripts --exclude Brewfile --exclude VENDORED_DO_NOT_MERGE.md /tmp/cesgo/services/evaluations/generated/go/ core/third_party/cesgo/
cd core && go mod tidy && go mod vendor
```
