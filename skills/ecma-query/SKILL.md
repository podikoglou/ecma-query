---
name: ecma-query
description: Query the ECMAScript specification from the CLI. Use when you need to look up spec sections, abstract operations, built-in methods, grammar productions, or cross-references — e.g. "what does ToNumber do?", "show me the Promise.all algorithm", "find the spec for Array.prototype.map", "what calls IsCallable?", or "show me the grammar for ArrowFunction".
---

# ecma-query

CLI for looking up the ECMAScript spec. Always use `--format md` for readable markdown output (default is JSON).

Bare `ecma-query <id>` is shorthand for `ecma-query get <id>`.

## Commands

### `get <id>` — exact entity lookup

Resolution order: section number (`7.1.4`) → anchor ID (`sec-toprimitive`) → abstract operation (`ToNumber`) → built-in method (`Array.prototype.map`) → internal slot (`[[Prototype]]`) → spec type (`PropertyDescriptor`). Operation names are case-insensitive and normalized (spaces/parens stripped).

```
ecma-query get ToNumber --format md
ecma-query get Array.prototype.map --format md
ecma-query get 7.1.4 --format md
ecma-query get sec-toprimitive --format md
```

| Flag | Default | What it does |
|------|---------|--------------|
| `--max-tokens N` | 0 (off) | Truncate output at ~N tokens (heuristic: word count / 0.75) |
| `--brief` | false | **JSON only.** Signature + summary, no steps |
| `--depth N` | 1 | **JSON only.** `> 0` includes algorithm steps, `0` = no steps (= `--brief`) |

### `search <query>` — fuzzy find

Full-text search with ranked results and snippets. Tokens shorter than 3 chars and common stop words are silently excluded. Use when you don't know the exact name.

```
ecma-query search "promise" --format md
ecma-query search "promise" --kind operation --format md
ecma-query search "promise" --kind grammar --limit 5 --format md
```

| Flag | Default | What it does |
|------|---------|--------------|
| `--limit N` | 10 | Max results |
| `--kind` | (all) | Filter: `operation`, `method`, `section`, `type`, `grammar`, `slot` |

### `toc [section]` — spec tree navigation

Without args: top-level chapters. With a section number: children of that section.

```
ecma-query toc --format md
ecma-query toc 27.2 --format md
```

| Flag | Default | What it does |
|------|---------|--------------|
| `--depth N` | 2 | Nesting levels to show |

### `xref <id>` — cross-references

Shows incoming and outgoing references. Same resolution and normalization as `get`.

```
ecma-query xref ToNumber --format md
ecma-query xref IsCallable --format md --direction incoming
```

| Flag | Default | What it does |
|------|---------|--------------|
| `--direction` | `both` | `incoming` (who references this), `outgoing` (what this references), `both` |

### `grammar <ProductionName>` — grammar productions

Case-sensitive. Use `search --kind grammar` to discover production names. **This command has its own `--format` flag** (default `ebnf`) — the global `--format` is ignored.

```
ecma-query grammar Identifier              # defaults to EBNF
ecma-query grammar Identifier --format md
ecma-query grammar Identifier --format json
```

## Workflow

1. **Don't know the exact name?** → `search <query> --format md`
2. **Need the full spec?** → `get <id> --format md`
3. **Output too long?** → add `--max-tokens N`
4. **Need to see references?** → `xref <id> --format md`
5. **Need a grammar rule?** → `grammar <Name> --format md`
6. **Need to browse structure?** → `toc [section] --format md`

## Notes

- `--spec` is defined but not implemented — don't bother passing it.
- All output is deterministic, no ANSI codes.
- Operation names are normalized: `get tonumber` works the same as `get ToNumber`.
