
# Getting Started

## Installation

### install with go

```bash
go install github.com/dfirebaugh/bbook@latest
```
> Installing with go should work cross platform

### Download a release
[https://github.com/dfirebaugh/bbook/releases/latest](https://github.com/dfirebaugh/bbook/releases/latest)

### Building from source

1. Pull down the git repository.
2. cd into the directory
3. run `go build`

```bash
git clone github.com/dfirebaugh/bbook
cd bbook
go build
```

## Make a book

Navigate to a directory where you would like to create your markdown files.

Run the following command

`bbook init <book name>`

e.g.

```bash
bbook init bbook
```

This will create a `book.toml` file that will act as your config.
It will also create a `SUMMARY.md` file.  This is a special markdown file that will be parsed in as the navigation sidebar.

### Serve the static site locally

```bash
bbook serve
```

`bbook serve` will serve the `.book` dir by default.  It will also create any markdown files that are linked in the `SUMMARY.md` file (if they don't exist).

> Any file changes will be automatically be rebuilt (however, auto page reload isn't working yet).

