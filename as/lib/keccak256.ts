// Pure AssemblyScript Keccak-256 implementation (Ethereum variant)
// Keccak-256 uses padding byte 0x01, NOT SHA3-256 which uses 0x06.

// Keccak-f[1600] round constants
const RC: u64[] = [
  0x0000000000000001, 0x0000000000008082, 0x800000000000808a, 0x8000000080008000,
  0x000000000000808b, 0x0000000080000001, 0x8000000080008081, 0x8000000000008009,
  0x000000000000008a, 0x0000000000000088, 0x0000000080008009, 0x000000008000000a,
  0x000000008000808b, 0x800000000000008b, 0x8000000000008089, 0x8000000000008003,
  0x8000000000008002, 0x8000000000000080, 0x000000000000800a, 0x800000008000000a,
  0x8000000080008081, 0x8000000000008080, 0x0000000080000001, 0x8000000080008008,
];

// Rotation offsets for rho step
const ROTATIONS: i32[] = [
   0,  1, 62, 28, 27,
  36, 44,  6, 55, 20,
   3, 10, 43, 25, 39,
  41, 45, 15, 21,  8,
  18,  2, 61, 56, 14,
];

// Pi step permutation indices
const PI: i32[] = [
   0, 10, 20,  5, 15,
  16,  1, 11, 21,  6,
   7, 17,  2, 12, 22,
  23,  8, 18,  3, 13,
  14, 24,  9, 19,  4,
];

const RATE: i32 = 136; // bytes (1088 bits) for Keccak-256

/** Keccak-f[1600] permutation: 24 rounds of theta/rho/pi/chi/iota */
function keccakF(state: StaticArray<u64>): void {
  for (let round = 0; round < 24; round++) {
    // Theta
    const c0 = state[0] ^ state[5] ^ state[10] ^ state[15] ^ state[20];
    const c1 = state[1] ^ state[6] ^ state[11] ^ state[16] ^ state[21];
    const c2 = state[2] ^ state[7] ^ state[12] ^ state[17] ^ state[22];
    const c3 = state[3] ^ state[8] ^ state[13] ^ state[18] ^ state[23];
    const c4 = state[4] ^ state[9] ^ state[14] ^ state[19] ^ state[24];

    const d0 = c4 ^ rotl<u64>(c1, 1);
    const d1 = c0 ^ rotl<u64>(c2, 1);
    const d2 = c1 ^ rotl<u64>(c3, 1);
    const d3 = c2 ^ rotl<u64>(c4, 1);
    const d4 = c3 ^ rotl<u64>(c0, 1);

    state[ 0] ^= d0; state[ 5] ^= d0; state[10] ^= d0; state[15] ^= d0; state[20] ^= d0;
    state[ 1] ^= d1; state[ 6] ^= d1; state[11] ^= d1; state[16] ^= d1; state[21] ^= d1;
    state[ 2] ^= d2; state[ 7] ^= d2; state[12] ^= d2; state[17] ^= d2; state[22] ^= d2;
    state[ 3] ^= d3; state[ 8] ^= d3; state[13] ^= d3; state[18] ^= d3; state[23] ^= d3;
    state[ 4] ^= d4; state[ 9] ^= d4; state[14] ^= d4; state[19] ^= d4; state[24] ^= d4;

    // Rho + Pi (combined)
    const temp = new StaticArray<u64>(25);
    for (let i = 0; i < 25; i++) {
      unchecked(temp[PI[i]] = rotl<u64>(state[i], ROTATIONS[i]));
    }

    // Chi
    for (let y = 0; y < 5; y++) {
      const base = y * 5;
      const t0 = unchecked(temp[base + 0]);
      const t1 = unchecked(temp[base + 1]);
      const t2 = unchecked(temp[base + 2]);
      const t3 = unchecked(temp[base + 3]);
      const t4 = unchecked(temp[base + 4]);
      unchecked(state[base + 0] = t0 ^ (~t1 & t2));
      unchecked(state[base + 1] = t1 ^ (~t2 & t3));
      unchecked(state[base + 2] = t2 ^ (~t3 & t4));
      unchecked(state[base + 3] = t3 ^ (~t4 & t0));
      unchecked(state[base + 4] = t4 ^ (~t0 & t1));
    }

    // Iota
    unchecked(state[0] ^= RC[round]);
  }
}

/**
 * Compute Keccak-256 hash of input data.
 * Returns 32 bytes.
 */
export function keccak256(data: Uint8Array): Uint8Array {
  const state = new StaticArray<u64>(25);

  // Absorb phase
  let offset = 0;
  while (offset + RATE <= data.length) {
    // XOR a full rate block into the state
    for (let i = 0; i < RATE; i += 8) {
      const laneIdx = i >> 3;
      const b0 = <u64>data[offset + i];
      const b1 = <u64>data[offset + i + 1] << 8;
      const b2 = <u64>data[offset + i + 2] << 16;
      const b3 = <u64>data[offset + i + 3] << 24;
      const b4 = <u64>data[offset + i + 4] << 32;
      const b5 = <u64>data[offset + i + 5] << 40;
      const b6 = <u64>data[offset + i + 6] << 48;
      const b7 = <u64>data[offset + i + 7] << 56;
      unchecked(state[laneIdx] ^= b0 | b1 | b2 | b3 | b4 | b5 | b6 | b7);
    }
    keccakF(state);
    offset += RATE;
  }

  // Pad: remaining bytes + Keccak padding
  const remaining = data.length - offset;
  const padBlock = new Uint8Array(RATE);
  for (let i = 0; i < remaining; i++) {
    padBlock[i] = data[offset + i];
  }
  // Keccak padding: 0x01 after data, 0x80 at end of block
  padBlock[remaining] = 0x01;
  padBlock[RATE - 1] |= 0x80;

  // XOR pad block into state
  for (let i = 0; i < RATE; i += 8) {
    const laneIdx = i >> 3;
    const b0 = <u64>padBlock[i];
    const b1 = <u64>padBlock[i + 1] << 8;
    const b2 = <u64>padBlock[i + 2] << 16;
    const b3 = <u64>padBlock[i + 3] << 24;
    const b4 = <u64>padBlock[i + 4] << 32;
    const b5 = <u64>padBlock[i + 5] << 40;
    const b6 = <u64>padBlock[i + 6] << 48;
    const b7 = <u64>padBlock[i + 7] << 56;
    unchecked(state[laneIdx] ^= b0 | b1 | b2 | b3 | b4 | b5 | b6 | b7);
  }
  keccakF(state);

  // Squeeze: extract first 32 bytes (4 lanes)
  const hash = new Uint8Array(32);
  for (let lane = 0; lane < 4; lane++) {
    const val = unchecked(state[lane]);
    const base = lane * 8;
    hash[base + 0] = <u8>(val & 0xff);
    hash[base + 1] = <u8>((val >> 8) & 0xff);
    hash[base + 2] = <u8>((val >> 16) & 0xff);
    hash[base + 3] = <u8>((val >> 24) & 0xff);
    hash[base + 4] = <u8>((val >> 32) & 0xff);
    hash[base + 5] = <u8>((val >> 40) & 0xff);
    hash[base + 6] = <u8>((val >> 48) & 0xff);
    hash[base + 7] = <u8>((val >> 56) & 0xff);
  }

  return hash;
}
