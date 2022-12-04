#!/bin/bash
set -ex

if ! command -v github-release >/dev/null; then
  echo must install github-release
  echo go get github.com/github-release/github-release
  exit 1
fi

if [ -z "$GITHUB_TOKEN" ]; then
  echo must set GITHUB_TOKEN >&2
  exit 1
fi

TAG=0.2.1

git tag -a v$TAG -m "release v$TAG"

git push origin master --tags

./compile

github-release release \
  --user acarl005 \
  --repo ls-go \
  --tag v$TAG \
  --name v$TAG


github-release upload \
  --user acarl005 \
  --repo ls-go \
  --tag v$TAG \
  --name ls-go-darwin-amd64 \
  --file ls-go-darwin-amd64

github-release upload \
  --user acarl005 \
  --repo ls-go \
  --tag v$TAG \
  --name ls-go-linux-amd64 \
  --file ls-go-linux-amd64

github-release upload \
  --user acarl005 \
  --repo ls-go \
  --tag v$TAG \
  --name ls-go-linux-arm64 \
  --file ls-go-linux-arm64

github-release upload \
  --user acarl005 \
  --repo ls-go \
  --tag v$TAG \
  --name ls-go-linux-386 \
  --file ls-go-linux-386
