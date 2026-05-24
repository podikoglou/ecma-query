# ecma-query

CLI for querying the ECMAScript specification. Designed for LLM agent consumption.

## Install

```
go install github.com/podikoglou/ecma-query@latest
```

## Usage

```
ecma-query [--format json|md] [--spec latest] <command> [args...]
```

Shorthand — no subcommand defaults to `get`:

```
ecma-query ToNumber    # same as: ecma-query get ToNumber
```

### Commands

**get** — retrieve a spec entity by exact identifier:
```
ecma-query get ToNumber
ecma-query get 7.1.4
ecma-query get sec-toprimitive
ecma-query get Array.prototype.map
```

Options: `--depth <n>`, `--brief`, `--steps-only`, `--max-tokens <n>`, `--chunk <n>`

**search** — full-text search with ranked snippets:
```
ecma-query search "promise resolution" --limit 5
ecma-query search "abstract operation" --kind operation
```

**toc** — table of contents / structural outline:
```
ecma-query toc              # top-level
ecma-query toc 23.1 --depth 3   # children of 23.1
```

**xref** — cross-reference resolution:
```
ecma-query xref OrdinaryObjectCreate --direction incoming
```

**grammar** — grammar production lookup (defaults to EBNF):
```
ecma-query grammar ArrowFunction
ecma-query grammar IfStatement --format json
```

### Output formats

- `json` (default) — structured JSON to stdout
- `md` — markdown with headings, breadcrumbs, and links
- `ebnf` — raw grammar text (grammar command only)

### Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Not found (structured error on stdout) |
| 2 | Ambiguous match (disambiguation list on stdout) |

### Build from source

```
git clone https://github.com/podikoglou/ecma-query.git
cd ecma-query
make
```

## License

MIT
