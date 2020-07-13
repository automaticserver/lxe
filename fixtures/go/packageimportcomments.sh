#!/bin/bash

# see: https://gist.github.com/dionysius/c803235e2e94353ed72be99ef208d428

set -o nounset -e

# we want to exit if current working dir is no go project (would error and thus exit with set -e)
go list -m >/dev/null

while read -r file; do
  reldir=$(dirname "$file")
  base=$(basename "$file")

  # exclude go module files
  if [[ "$base" == "go.mod" ]] || [[ "$base" == "go.sum" ]]; then
    continue
  fi

  package=$(go list -f '{{ .Name }}' ./"$reldir" 2>/dev/null || true)
  
  # if not a package due to build constraints
  if [[ "$package" == "" ]]; then
    continue
  fi

  # a main package is not importable
  if [[ "$package" == "main" ]]; then
    continue
  fi

  # an internal package is not importable
  if [[ "$package" =~ (^|/)"internal"(/|$) ]]; then
    continue
  fi

  importpath=$(go list -f '{{ .ImportPath }}' ./"$reldir")


  sed -i "s|^package $package.*|package $package // import \"$importpath\"|" "$file"
done < "${1:-/dev/stdin}"
