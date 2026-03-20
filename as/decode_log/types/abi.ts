/** A single input parameter from an ABI event definition. */
export class AbiInput {
  name: string;
  typeName: string;
  indexed: bool;

  constructor(name: string, typeName: string, indexed: bool) {
    this.name = name;
    this.typeName = typeName;
    this.indexed = indexed;
  }
}

/** A parsed ABI event with its precomputed selector hash. */
export class AbiEvent {
  name: string;
  signature: string;
  selector: string;
  inputs: AbiInput[];

  constructor(name: string, signature: string, selector: string, inputs: AbiInput[]) {
    this.name = name;
    this.signature = signature;
    this.selector = selector;
    this.inputs = inputs;
  }
}
