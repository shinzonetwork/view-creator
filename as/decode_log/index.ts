import { JSON } from "assemblyscript-json/assembly";
import {
  TransportMessage,
  fromTransportVec,
  toTransportVec,
  nilPtr,
  JSON_TYPE_ID,
  EOS_TYPE_ID,
  ERROR_TYPE_ID,
} from "../lib/transport";
import { Parameters } from "./types/params";
import { InputLog } from "./types/input";
import { OutputLog } from "./types/output";
import { parseAbi, findMatchingEvent, decodeEvent } from "./decoder";

// @ts-ignore: decorator
@external("lens", "next")
declare function next(): usize;

// @ts-ignore: decorator
function abort(
  message: string | null,
  fileName: string | null,
  lineNumber: u32,
  columnNumber: u32
): void {
  unreachable();
}

let params: Parameters | null = null;

export function alloc(size: usize): usize {
  return heap.alloc(size);
}

export function set_param(ptr: usize): usize {
  const msg = fromTransportVec(ptr);
  if (msg.isEndOfStream) {
    return toTransportVec(ERROR_TYPE_ID, "unexpected end of stream in set_param");
  }

  const parsed = <JSON.Obj>JSON.parse(msg.payload);
  params = Parameters.fromJson(parsed);

  if (params!.abi.length == 0) {
    return toTransportVec(ERROR_TYPE_ID, "missing 'abi' parameter");
  }

  heap.free(ptr);
  return nilPtr();
}

export function transform(): usize {
  const ptr = next();
  const msg = fromTransportVec(ptr);

  if (msg.isEndOfStream) {
    heap.free(ptr);
    return toTransportVec(EOS_TYPE_ID, "");
  }

  if (msg.payload.length == 0) {
    heap.free(ptr);
    return toTransportVec(JSON_TYPE_ID, "{}");
  }

  const doc = <JSON.Obj>JSON.parse(msg.payload);
  const input = InputLog.fromJson(doc);
  const output = processLog(input);
  const json = output.toJson();

  heap.free(ptr);
  return toTransportVec(JSON_TYPE_ID, json);
}

function processLog(input: InputLog): OutputLog {
  const out = new OutputLog();
  out.hash = input.transaction.hash;
  out.from = input.transaction.from;
  out.to = input.transaction.to;
  out.blockNumber = input.blockNumber;
  out.logAddress = input.address;

  if (input.topics.length > 0 && params != null && params!.abi.length > 0) {
    const events = parseAbi(params!.abi);
    const matched = findMatchingEvent(events, input.topics[0]);
    if (matched != null) {
      out.event = matched.name;
      out.signature = matched.signature;
      out.arguments = decodeEvent(matched, input.topics, input.data);
    }
  }

  return out;
}
