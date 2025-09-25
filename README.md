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
./viewkit view deploy testdeploy --target devnet
```

**What each step does**

- `view init` create a versioned bundle on disk with metadata
- `view inspect` print bundle contents (queries, SDL, lenses, versions)
- `view add query` define the raw shape you plan to ingest (e.g. EVM logs)
- `view add sdl` define your GraphQL type; toggle persistence with `@materialized(if: true|false)`
- `view add lens` chainable WASM transforms that pre-process data
- `wallet generate` create a local signing key (store securely)
- `view deploy` publish the bundle to a target network
