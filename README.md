# View Creator (Viewkit)

**Viewkit helps you initialize, manage, and publish Shinzo views through a simple CLI interface.**

A **view** is a versioned bundle that can include:
- **Queries** the raw data shape you want to ingest
- **SDL** a GraphQL schema that describes how data is modeled or materialized
- **Lenses** WebAssembly transforms for filtering, decoding, or reshaping data
- **Wallet** credentials used to sign deployments to a target network

This guide gets you from zero to a deployed view on `devnet` in minutes.

---

## Install

```bash
# from the repo root
make build
```

You can then run `viewkit` via `./build/viewkit` or simply `viewkit` if you exported PATH.

---

## Wasmer runtime (macOS, Apple Silicon)

### Why
Viewkit can execute WebAssembly **lenses** locally for validation and preview. The Go package `wasmer-go` uses a native dynamic library (`libwasmer.dylib`) at runtime. If the loader cannot find it, commands that touch lenses will fail with an error like “image not found.”

### What the env vars do
- `WASMER_LIB_PATH` tells `wasmer-go` where the Wasmer dynamic library lives.
- `DYLD_LIBRARY_PATH` tells the macOS dynamic loader where to search for `.dylib` files when a process starts. We prepend the Wasmer directory so the loader can find `libwasmer.dylib`.

### Install the module
You already have the module via `go get`, but here’s the explicit call if needed:
```bash
go get github.com/wasmerio/wasmer-go@v1.0.4
```

### Add the env vars to your shell

```bash
# Wasmer Go native libs (Apple Silicon)
export WASMER_ROOT="$(go env GOPATH)/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/darwin-aarch64"
export WASMER_LIB_PATH="$WASMER_ROOT"
export DYLD_LIBRARY_PATH="$WASMER_ROOT:$DYLD_LIBRARY_PATH"
```

Apply it now:
```bash
source ~/.zshrc
```

One-liners to append automatically:
```bash
echo 'export WASMER_ROOT="$(go env GOPATH)/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/darwin-aarch64"' >> ~/.zshrc
echo 'export WASMER_LIB_PATH="$WASMER_ROOT"' >> ~/.zshrc
echo 'export DYLD_LIBRARY_PATH="$WASMER_ROOT:$DYLD_LIBRARY_PATH"' >> ~/.zshrc
source ~/.zshrc
```

---

## Quick start

Create, inspect, enrich, and deploy a view named **testdeploy**:

```bash
# 1) initialize the view bundle
./viewkit view init testdeploy

# 2) inspect the bundle
./viewkit view inspect testdeploy

# 3) add a query (raw event shape to ingest)
./viewkit view add query "Log {address topics data transactionHash blockNumber}" --name testdeploy

# 4) add SDL (how data is modeled / stored)
./viewkit view add sdl "type FilteredAndDecodedLogs @materialized(if: false) {transactionHash: String}" --name testdeploy

# 5) attach a lens (WASM transform to filter by address)
./viewkit view add lens   --args '{"src":"address", "value":"0x1e3aA9fE4Ef01D3cB3189c129a49E3C03126C636"}'   --label "filter"   --url "https://raw.githubusercontent.com/shinzonetwork/wasm-bucket/main/bucket/filter_transaction/filter_transaction.wasm"   --name testdeploy

# 6) create a wallet for deployments
./viewkit wallet generate

# 7) deploy to devnet
./viewkit view deploy testdeploy --target devnet --rpc http://34.29.171.79:8545/
```

**What each step does**

- `view init` create a versioned bundle on disk with metadata
- `view inspect` print bundle contents (queries, SDL, lenses, versions)
- `view add query` define the raw shape you plan to ingest (e.g. EVM logs)
- `view add sdl` define your GraphQL type; toggle persistence with `@materialized(if: true|false)`
- `view add lens` chainable WASM transforms that pre-process data
- `wallet generate` create a local signing key (store securely)
- `view deploy` publish the bundle to a target network
