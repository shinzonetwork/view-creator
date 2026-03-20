import { JSONEncoder } from "assemblyscript-json/assembly";

/** A single decoded event argument. */
export class DecodedArgument {
  name: string;
  type: string;
  value: string;

  constructor(name: string, type: string, value: string) {
    this.name = name;
    this.type = type;
    this.value = value;
  }
}

/** Decoded event log output document. */
export class OutputLog {
  hash: string = "";
  from: string = "";
  to: string = "";
  blockNumber: i64 = 0;
  logAddress: string = "";
  event: string = "Unknown";
  signature: string = "";
  arguments: DecodedArgument[] = [];

  /** Serialize this output to a JSON string using JSONEncoder. */
  toJson(): string {
    const enc = new JSONEncoder();
    enc.pushObject(null);

    enc.setString("hash", this.hash);
    enc.setInteger("blockNumber", this.blockNumber);
    enc.setString("from", this.from);
    enc.setString("to", this.to);
    enc.setString("logAddress", this.logAddress);
    enc.setString("event", this.event);
    enc.setString("signature", this.signature);

    enc.pushArray("arguments");
    for (let i = 0; i < this.arguments.length; i++) {
      const arg = this.arguments[i];
      enc.pushObject(null);
      enc.setString("name", arg.name);
      enc.setString("type", arg.type);
      enc.setString("value", arg.value);
      enc.popObject();
    }
    enc.popArray();

    enc.popObject();
    return enc.toString();
  }
}
