# Temporary copy of the CES Go SDK: do not merge

> [!WARNING]
> TODO: Delete this directory and the `replace` directive in `core/go.mod`
> once `github.com/wandb/cwevals-go` publishes a release, and require that
> version instead. This copy exists only for early development and must not
> be merged.

This is a copy of the generated CoreWeave Evaluations (CES) Go SDK from
`services/evaluations/generated/go` in wandb/core at commit
`35fbeaaf874fcf2ac55727e4474843eb254a7c29` (`kelu/ces-stainless-case`), without
the files that `wandb/cwevals-go` releases leave out or that only run its
own tooling.

To refresh it from a wandb/core checkout at `$CORE`:

```bash
git -C "$CORE" archive <commit> services/evaluations/generated/go | tar -x -C /tmp/cwevals-go
rsync -a --delete --exclude examples --exclude scripts --exclude lib --exclude Brewfile --exclude .devcontainer --exclude .github --exclude .stats.yml --exclude .gitignore --exclude VENDORED_DO_NOT_MERGE.md /tmp/cwevals-go/services/evaluations/generated/go/ core/third_party/cwevals-go/
cd core && go mod tidy && go mod vendor
```
