// Hex encoding/decoding utilities for Ethereum data

const HEX_CHARS: string = "0123456789abcdef";

/** Convert a hex character code to its nibble value (0-15). Returns 0 for invalid chars. */
function hexCharToNibble(c: i32): u8 {
  if (c >= 48 && c <= 57) return <u8>(c - 48);        // '0'-'9'
  if (c >= 97 && c <= 102) return <u8>(c - 97 + 10);  // 'a'-'f'
  if (c >= 65 && c <= 70) return <u8>(c - 65 + 10);   // 'A'-'F'
  return 0;
}

/** Decode a hex string to bytes. Strips leading "0x" if present. */
export function hexDecode(hex: string): Uint8Array {
  let start = 0;
  if (hex.length >= 2 && hex.charAt(0) == "0" && (hex.charAt(1) == "x" || hex.charAt(1) == "X")) {
    start = 2;
  }

  const cleanLen = hex.length - start;
  // Handle odd-length hex strings by padding with a leading zero
  const byteLen = (cleanLen + 1) >> 1;
  const result = new Uint8Array(byteLen);

  let offset = 0;
  let i = start;
  if (cleanLen & 1) {
    // Odd length: first byte is just a single nibble
    result[0] = hexCharToNibble(hex.charCodeAt(i));
    i++;
    offset = 1;
  }

  while (i < hex.length) {
    const hi = hexCharToNibble(hex.charCodeAt(i));
    const lo = hexCharToNibble(hex.charCodeAt(i + 1));
    result[offset] = (hi << 4) | lo;
    i += 2;
    offset++;
  }

  return result;
}

/** Encode bytes to a lowercase hex string (no "0x" prefix). */
export function hexEncode(bytes: Uint8Array): string {
  let result = "";
  for (let i = 0; i < bytes.length; i++) {
    const b = bytes[i];
    result += HEX_CHARS.charAt((b >> 4) & 0x0f);
    result += HEX_CHARS.charAt(b & 0x0f);
  }
  return result;
}

/**
 * Parse a uint256 from a hex string (up to 128 bits supported via two u64s).
 * Returns the decimal string representation.
 * For values > u64 max, we parse into high and low 64-bit parts.
 */
export function hexToDecimalString(hex: string): string {
  let clean = hex;
  if (clean.length >= 2 && clean.charAt(0) == "0" && (clean.charAt(1) == "x" || clean.charAt(1) == "X")) {
    clean = clean.substring(2);
  }

  // Strip leading zeros
  let startIdx = 0;
  while (startIdx < clean.length - 1 && clean.charAt(startIdx) == "0") {
    startIdx++;
  }
  clean = clean.substring(startIdx);

  // If it fits in a u64 (16 hex chars = 64 bits), parse directly
  if (clean.length <= 16) {
    let value: u64 = 0;
    for (let i = 0; i < clean.length; i++) {
      value = (value << 4) | <u64>hexCharToNibble(clean.charCodeAt(i));
    }
    return value.toString();
  }

  // For larger values, use digit-by-digit decimal arithmetic
  return hexToDecimalBignum(clean);
}

/** Convert arbitrarily large hex to decimal string using digit-by-digit multiplication. */
function hexToDecimalBignum(hex: string): string {
  // Represent the decimal result as an array of digits (least significant first)
  let digits: u8[] = [0];

  for (let i = 0; i < hex.length; i++) {
    const nibble = hexCharToNibble(hex.charCodeAt(i));

    // Multiply current result by 16
    let carry: u32 = 0;
    for (let j = 0; j < digits.length; j++) {
      const prod = <u32>digits[j] * 16 + carry;
      digits[j] = <u8>(prod % 10);
      carry = prod / 10;
    }
    while (carry > 0) {
      digits.push(<u8>(carry % 10));
      carry = carry / 10;
    }

    // Add the nibble
    carry = <u32>nibble;
    for (let j = 0; j < digits.length && carry > 0; j++) {
      const sum = <u32>digits[j] + carry;
      digits[j] = <u8>(sum % 10);
      carry = sum / 10;
    }
    while (carry > 0) {
      digits.push(<u8>(carry % 10));
      carry = carry / 10;
    }
  }

  // Convert digits array to string (reverse order — most significant first)
  let result = "";
  for (let i = digits.length - 1; i >= 0; i--) {
    result += String.fromCharCode(48 + <i32>digits[i]);
  }
  return result;
}
