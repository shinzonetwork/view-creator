import { JSON } from "assemblyscript-json/assembly";

/** Lens parameters containing the ABI definition string. */
export class Parameters {
  abi: string = "";

  /**
   * Parse Parameters from a JSON.Obj.
   * Accepts "abi" as either a JSON string or a JSON array:
   *   {"abi": "[{...}]"}       ← string (from Go map[string]any with string value)
   *   {"abi": [{...}]}         ← array  (from CLI --args with inline JSON array)
   */
  static fromJson(obj: JSON.Obj): Parameters {
    const p = new Parameters();

    // Try string first (e.g. from Go service layer or escaped CLI input)
    const abiStr = obj.getString("abi");
    if (abiStr != null) {
      p.abi = abiStr.valueOf();
      return p;
    }

    // Fall back to array — stringify it so the rest of the pipeline works
    const abiArr = obj.getArr("abi");
    if (abiArr != null) {
      p.abi = abiArr.stringify();
      return p;
    }

    return p;
  }
}
