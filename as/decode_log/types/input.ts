import { JSON } from "assemblyscript-json/assembly";

/** Nested transaction context from the input log document. */
export class Transaction {
  hash: string = "";
  from: string = "";
  to: string = "";
}

/** Raw Ethereum log document as received from the LensVM pipeline. */
export class InputLog {
  address: string = "";
  topics: string[] = [];
  data: string = "";
  blockNumber: i64 = 0;
  transaction: Transaction = new Transaction();

  /** Parse an InputLog from a JSON.Obj, encapsulating all nullable-check extraction. */
  static fromJson(doc: JSON.Obj): InputLog {
    const log = new InputLog();

    // address
    const addrVal = doc.getString("address");
    if (addrVal != null) log.address = addrVal.valueOf();

    // blockNumber
    const blockNumVal = doc.getInteger("blockNumber");
    if (blockNumVal != null) log.blockNumber = blockNumVal.valueOf();

    // data
    const dataVal = doc.getString("data");
    if (dataVal != null) log.data = dataVal.valueOf();

    // topics
    const topicsArr = doc.getArr("topics");
    if (topicsArr != null) {
      const values = topicsArr.valueOf();
      for (let i = 0; i < values.length; i++) {
        if (values[i].isString) {
          log.topics.push((<JSON.Str>values[i]).valueOf());
        }
      }
    }

    // transaction (nested object)
    const txObj = doc.getObj("transaction");
    if (txObj != null) {
      const h = txObj.getString("hash");
      if (h != null) log.transaction.hash = h.valueOf();
      const f = txObj.getString("from");
      if (f != null) log.transaction.from = f.valueOf();
      const t = txObj.getString("to");
      if (t != null) log.transaction.to = t.valueOf();
    }

    return log;
  }
}
