import { readFileSync } from "fs";
import { join } from "path";

export const JSON_TYPE_ID = 1;
export const EOS_TYPE_ID = 127;

export type WasmExports = {
  memory: WebAssembly.Memory;
  alloc: (size: number) => number;
  set_param: (ptr: number) => number;
  transform: () => number;
};

/** Encode a JSON string into a LensVM transport vector. */
export function encodeTransport(typeId: number, payload: string): Uint8Array {
  const payloadBytes = new TextEncoder().encode(payload);
  const buf = new Uint8Array(1 + 4 + payloadBytes.length);
  const view = new DataView(buf.buffer);
  view.setInt8(0, typeId);
  view.setUint32(1, payloadBytes.length, true); // little-endian
  buf.set(payloadBytes, 5);
  return buf;
}

/** Decode a LensVM transport vector from wasm memory. */
export function decodeTransport(memory: WebAssembly.Memory, ptr: number): { typeId: number; payload: string } {
  const view = new DataView(memory.buffer);
  const typeId = view.getInt8(ptr);
  if (typeId === EOS_TYPE_ID) {
    return { typeId, payload: "" };
  }
  const len = view.getUint32(ptr + 1, true);
  const payloadBytes = new Uint8Array(memory.buffer, ptr + 5, len);
  const payload = new TextDecoder().decode(payloadBytes);
  return { typeId, payload };
}

/** Write a transport vector into wasm memory and return the pointer. */
export function writeToMemory(
  exports: { memory: WebAssembly.Memory; alloc: (size: number) => number },
  data: Uint8Array
): number {
  const ptr = exports.alloc(data.length);
  new Uint8Array(exports.memory.buffer).set(data, ptr);
  return ptr;
}

/** Queue of transport vectors to be returned by the mocked next() function. */
export let nextQueue: Uint8Array[] = [];

export function setNextQueue(queue: Uint8Array[]) {
  nextQueue = queue;
}

export async function loadWasm(wasmDir: string): Promise<WasmExports> {
  const wasmPath = join(wasmDir, "..", "build", "decode_log", "decode_log.wasm");
  const wasmBuffer = readFileSync(wasmPath);

  let memoryRef: WebAssembly.Memory;
  let allocFn: (size: number) => number;

  const { instance } = await WebAssembly.instantiate(wasmBuffer, {
    lens: {
      next(): number {
        if (nextQueue.length === 0) {
          const eos = encodeTransport(EOS_TYPE_ID, "");
          const ptr = allocFn(eos.length);
          new Uint8Array(memoryRef.buffer).set(eos, ptr);
          return ptr;
        }
        const vec = nextQueue.shift()!;
        const ptr = allocFn(vec.length);
        new Uint8Array(memoryRef.buffer).set(vec, ptr);
        return ptr;
      },
    },
    env: {
      abort: () => { throw new Error("abort called"); },
    },
  });

  const exports = instance.exports as WasmExports;

  memoryRef = exports.memory;
  allocFn = exports.alloc;

  return exports;
}
