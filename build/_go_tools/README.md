# Go binary tools dependencies

Here stored go binary dependencies and tools to manage them.

* ls-imports - directory
* tools - directory with binary dependencies described in go.mod file and vendor dir with all required for dependencies

## How to manage dependencies

add dependency

update dependency

remove dependency

## Good to know

go tools doesn't work well in directories with underscore in its names (i.e. `go vendor`), so we forced to use subdirectory `tools/` here.
