import { describe, test, expect } from "bun:test";
import {
  JSON_TYPE_ID, EOS_TYPE_ID,
  encodeTransport, decodeTransport, writeToMemory,
  loadWasm, setNextQueue,
} from "../lib/test-helpers";
import { ERC20_ABI, SAMPLE_LOG } from "../lib/test-data";

describe("decode_log wasm", () => {
  test("exports alloc, set_param, transform and imports lens.next", async () => {
    const exports = await loadWasm(import.meta.dir);
    expect(typeof exports.alloc).toBe("function");
    expect(typeof exports.set_param).toBe("function");
    expect(typeof exports.transform).toBe("function");
    expect(exports.memory).toBeInstanceOf(WebAssembly.Memory);
  });

  test("set_param accepts ABI and returns nil", async () => {
    const exports = await loadWasm(import.meta.dir);
    const paramJson = JSON.stringify({ abi: ERC20_ABI });
    const vec = encodeTransport(JSON_TYPE_ID, paramJson);
    const ptr = writeToMemory(exports, vec);

    const resultPtr = exports.set_param(ptr);
    const view = new DataView(exports.memory.buffer);
    const typeId = view.getInt8(resultPtr);
    expect(typeId).toBe(0); // NIL_TYPE_ID = success
  });

  test("transform decodes an ERC-20 Transfer event", async () => {
    const exports = await loadWasm(import.meta.dir);

    // Set ABI parameter
    const paramJson = JSON.stringify({ abi: ERC20_ABI });
    const paramVec = encodeTransport(JSON_TYPE_ID, paramJson);
    const paramPtr = writeToMemory(exports, paramVec);
    exports.set_param(paramPtr);

    // Queue the log document for next()
    setNextQueue([encodeTransport(JSON_TYPE_ID, SAMPLE_LOG)]);

    const resultPtr = exports.transform();
    const result = decodeTransport(exports.memory, resultPtr);

    expect(result.typeId).toBe(JSON_TYPE_ID);

    const output = JSON.parse(result.payload);
    expect(output.hash).toBe("0xabc123");
    expect(output.from).toBe("0xsender");
    expect(output.to).toBe("0xreceiver");
    expect(output.blockNumber).toBe(18500000);
    expect(output.logAddress).toBe("0xdac17f958d2ee523a2206206994597c13d831ec7");
    expect(output.event).toBe("Transfer");
    expect(output.signature).toBe("Transfer(address,address,uint256)");

    // Check arguments
    expect(output.arguments).toHaveLength(3);

    expect(output.arguments[0].name).toBe("from");
    expect(output.arguments[0].type).toBe("address");
    expect(output.arguments[0].value).toBe("0xab5801a7d398351b8be11c439e05c5b3259aec9b");

    expect(output.arguments[1].name).toBe("to");
    expect(output.arguments[1].type).toBe("address");
    expect(output.arguments[1].value).toBe("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2");

    expect(output.arguments[2].name).toBe("value");
    expect(output.arguments[2].type).toBe("uint256");
    expect(output.arguments[2].value).toBe("2000000000"); // 0x77359400 = 2000000000
  });

  test("transform returns EOS when input stream is empty", async () => {
    const exports = await loadWasm(import.meta.dir);

    // Set ABI
    const paramVec = encodeTransport(JSON_TYPE_ID, JSON.stringify({ abi: ERC20_ABI }));
    exports.set_param(writeToMemory(exports, paramVec));

    // Queue EOS
    setNextQueue([encodeTransport(EOS_TYPE_ID, "")]);

    const resultPtr = exports.transform();
    const result = decodeTransport(exports.memory, resultPtr);
    expect(result.typeId).toBe(EOS_TYPE_ID);
  });

  test("transform returns Unknown event when topic0 does not match ABI", async () => {
    const exports = await loadWasm(import.meta.dir);

    // Set ABI
    const paramVec = encodeTransport(JSON_TYPE_ID, JSON.stringify({ abi: ERC20_ABI }));
    exports.set_param(writeToMemory(exports, paramVec));

    const unknownLog = JSON.stringify({
      address: "0x1234",
      topics: ["0x0000000000000000000000000000000000000000000000000000000000000000"],
      data: "0x",
      blockNumber: 100,
      transaction: { hash: "0xdef", from: "0xa", to: "0xb" },
    });
    setNextQueue([encodeTransport(JSON_TYPE_ID, unknownLog)]);

    const resultPtr = exports.transform();
    const result = decodeTransport(exports.memory, resultPtr);

    expect(result.typeId).toBe(JSON_TYPE_ID);
    const output = JSON.parse(result.payload);
    expect(output.event).toBe("Unknown");
    expect(output.signature).toBe("");
    expect(output.arguments).toHaveLength(0);
  });
});
