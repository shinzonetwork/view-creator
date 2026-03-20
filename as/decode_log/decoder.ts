import { JSON } from "assemblyscript-json/assembly";
import { AbiEvent, AbiInput } from "./types/abi";
import { DecodedArgument } from "./types/output";
import { keccak256 } from "../lib/keccak256";
import { hexEncode, hexToDecimalString } from "../lib/hex";

/** Parse a JSON ABI string into typed AbiEvent objects. Only "event" types are extracted. */
export function parseAbi(abiJson: string): AbiEvent[] {
  const events: AbiEvent[] = [];
  const items = <JSON.Arr>JSON.parse(abiJson);
  const arr = items.valueOf();

  for (let i = 0; i < arr.length; i++) {
    const item = <JSON.Obj>arr[i];

    const typeField = item.getString("type");
    if (typeField == null || typeField.valueOf() != "event") continue;

    const nameField = item.getString("name");
    if (nameField == null) continue;
    const name = nameField.valueOf();

    const inputsArr = item.getArr("inputs");
    if (inputsArr == null) continue;

    const inputs: AbiInput[] = [];
    const inputValues = inputsArr.valueOf();
    const types: string[] = [];

    for (let j = 0; j < inputValues.length; j++) {
      const inp = <JSON.Obj>inputValues[j];

      const inpName = inp.getString("name");
      const inpType = inp.getString("type");
      if (inpName == null || inpType == null) continue;

      const inpIndexed = inp.getBool("indexed");
      const indexed = inpIndexed != null ? inpIndexed.valueOf() : false;

      inputs.push(new AbiInput(inpName.valueOf(), inpType.valueOf(), indexed));
      types.push(inpType.valueOf());
    }

    const sig = name + "(" + types.join(",") + ")";
    const selector = computeSelector(sig);
    events.push(new AbiEvent(name, sig, selector, inputs));
  }

  return events;
}

/** Compute the Keccak-256 selector hash for an event signature string. */
export function computeSelector(signature: string): string {
  const sigBytes = String.UTF8.encode(signature, false);
  const hash = keccak256(Uint8Array.wrap(sigBytes));
  return "0x" + hexEncode(hash);
}

/** Find the ABI event whose selector matches the given topic0 hash. */
export function findMatchingEvent(events: AbiEvent[], topic0: string): AbiEvent | null {
  const topic0Lower = topic0.toLowerCase();
  for (let i = 0; i < events.length; i++) {
    if (events[i].selector.toLowerCase() == topic0Lower) {
      return events[i];
    }
  }
  return null;
}

/** Decode an event's indexed and non-indexed parameters into typed DecodedArgument objects. */
export function decodeEvent(event: AbiEvent, topics: string[], data: string): DecodedArgument[] {
  const args: DecodedArgument[] = [];
  const nonIndexed: AbiInput[] = [];
  let topicIdx = 1;

  for (let i = 0; i < event.inputs.length; i++) {
    const inp = event.inputs[i];
    if (inp.indexed) {
      if (topicIdx < topics.length) {
        const value = decodeParam(inp.typeName, topics[topicIdx]);
        args.push(new DecodedArgument(inp.name, inp.typeName, value));
      }
      topicIdx++;
    } else {
      nonIndexed.push(inp);
    }
  }

  // Decode non-indexed parameters from data field
  if (data.length > 2) {
    let dataHex = data;
    if (dataHex.length >= 2 && dataHex.charAt(0) == "0" && (dataHex.charAt(1) == "x" || dataHex.charAt(1) == "X")) {
      dataHex = dataHex.substring(2);
    }

    let offset = 0;
    for (let i = 0; i < nonIndexed.length; i++) {
      if (offset + 64 > dataHex.length) break;
      const chunk = "0x" + dataHex.substring(offset, offset + 64);
      const value = decodeParam(nonIndexed[i].typeName, chunk);
      args.push(new DecodedArgument(nonIndexed[i].name, nonIndexed[i].typeName, value));
      offset += 64;
    }
  }

  return args;
}

/** Decode a single ABI-encoded parameter from its hex representation. */
export function decodeParam(typeName: string, hexData: string): string {
  let clean = hexData;
  if (clean.length >= 2 && clean.charAt(0) == "0" && (clean.charAt(1) == "x" || clean.charAt(1) == "X")) {
    clean = clean.substring(2);
  }

  if (
    typeName == "uint256" || typeName == "uint128" || typeName == "uint64" ||
    typeName == "uint32" || typeName == "uint16" || typeName == "uint8" ||
    typeName == "int256" || typeName == "int128" || typeName == "int64" ||
    typeName == "int32" || typeName == "int16" || typeName == "int8"
  ) {
    return hexToDecimalString(clean);
  }

  if (typeName == "address") {
    if (clean.length >= 40) {
      return "0x" + clean.substring(clean.length - 40);
    }
    return "0x" + clean;
  }

  if (typeName == "bool") {
    return clean.endsWith("1") ? "true" : "false";
  }

  if (typeName == "bytes32" || typeName == "bytes20" || typeName == "bytes16" || typeName == "bytes4") {
    return "0x" + clean;
  }

  return "unsupported type: " + typeName;
}
