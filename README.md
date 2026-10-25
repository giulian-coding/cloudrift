# cloudrift

A self hosted Setup for managing Container Runtimes and their assignments to users. This is a work in progress, but the goal is to provide a simple, self-contained API for managing users and their assignments to container runtimes.

## Preparation

```bash
setup.sh
```

## Run

Requires Go 1.27+.

```bash
go run ./api    # http://localhost:8080
```

## Plan

Namespaces are used to group users and their assignments. Each namespace has a name and a description.

---

Work in progress — not intended for production use.
