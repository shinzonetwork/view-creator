# Shinzo Lenses — AssemblyScript

AssemblyScript [LensVM](https://docs.source.network/lensvm) lenses for the [Shinzo](https://github.com/shinzonetwork) indexing infrastructure. Each lens is a WASM module that transforms blockchain data inside the LensVM runtime.

## Prerequisites

- [Bun](https://bun.sh/) >= 1.0
- Node.js >= 18 (used by AssemblyScript compiler)

## Setup

```bash
cd as
bun install
```

## Build

```bash
# Build all lenses
bun run build

# Build a specific lens
bun run build:decode_log

# Build with debug output (includes WAT + source maps)
bun run build:decode_log:debug
```

Output goes to `build/<lens_name>/`.

Then, a WASM file from the `build` directory can be used with a `viewkit` CLI to deplo

## Test

```bash
bun test
```

Tests live alongside their lens in each lens directory (e.g. `decode_log/decode_log.test.ts`).

## Adding a new lens

1. **Create a directory** for the lens:

   ```bash
   mkdir -p my_lens/types
   ```

2. **Add an entry point** at `my_lens/index.ts`. It must export the LensVM interface:

   ```ts
   import { fromTransportVec, toTransportVec, nilPtr, JSON_TYPE_ID, EOS_TYPE_ID, ERROR_TYPE_ID } from "../lib/transport";

   @external("lens", "next")
   declare function next(): usize;

   function abort(message: string | null, fileName: string | null, lineNumber: u32, columnNumber: u32): void {
     unreachable();
   }

   export function alloc(size: usize): usize {
     return heap.alloc(size);
   }

   export function set_param(ptr: usize): usize {
     // Parse params, return nilPtr() on success
     return nilPtr();
   }

   export function transform(): usize {
     const ptr = next();
     const msg = fromTransportVec(ptr);
     if (msg.isEndOfStream) {
       heap.free(ptr);
       return toTransportVec(EOS_TYPE_ID, "");
     }
     // Transform msg.payload and return result
     heap.free(ptr);
     return toTransportVec(JSON_TYPE_ID, result);
   }
   ```

3. **Add build targets** in `asconfig.json`:

   ```json
   {
     "targets": {
       "my_lens": {
         "outFile": "build/my_lens/my_lens.wasm",
         "textFile": "build/my_lens/my_lens.wat",
         "sourceMap": true,
         "optimizeLevel": 3,
         "shrinkLevel": 0,
         "use": ["abort=my_lens/index/abort"]
       }
     }
   }
   ```

4. **Add build scripts** in `package.json`:

   ```json
   {
     "scripts": {
       "build": "bun run build:decode_log && bun run build:my_lens",
       "build:my_lens": "asc my_lens/index.ts --target my_lens"
     }
   }
   ```

5. **Add tests** at `my_lens/my_lens.test.ts`.

## References

- [LensVM documentation](https://docs.source.network/lensvm)
- [AssemblyScript](https://www.assemblyscript.org/)
- [assemblyscript-json](https://github.com/aspect-build/assemblyscript-json)
- [Source Network](https://source.network/)
