// LensVM binary transport protocol
// Format: [TypeId: i8][Length: u32 LE][Payload: bytes]

export const JSON_TYPE_ID: i8 = 1;
export const EOS_TYPE_ID: i8 = 127;
export const ERROR_TYPE_ID: i8 = -1;
export const NIL_TYPE_ID: i8 = 0;

export class TransportMessage {
  typeId: i8;
  payload: string;

  constructor(typeId: i8, payload: string) {
    this.typeId = typeId;
    this.payload = payload;
  }

  get isEndOfStream(): bool {
    return this.typeId == EOS_TYPE_ID;
  }

  get isJson(): bool {
    return this.typeId == JSON_TYPE_ID;
  }

  get isError(): bool {
    return this.typeId == ERROR_TYPE_ID;
  }

  get isNil(): bool {
    return this.typeId == NIL_TYPE_ID;
  }
}

/** Parse incoming LensVM transport vector from a raw memory pointer. */
export function fromTransportVec(ptr: usize): TransportMessage {
  const typeId = load<i8>(ptr);
  if (typeId == EOS_TYPE_ID) {
    return new TransportMessage(EOS_TYPE_ID, "");
  }
  const len = load<u32>(ptr + 1);
  const payload = String.UTF8.decodeUnsafe(ptr + 1 + 4, len, false);
  return new TransportMessage(typeId, payload);
}

/** Write an outgoing LensVM transport vector and return a pointer to it. */
export function toTransportVec(typeId: i8, message: string): usize {
  const len = String.UTF8.byteLength(message, false);
  const total = 1 + 4 + len;
  const ptr = heap.alloc(total);

  store<i8>(ptr, typeId);

  if (typeId == EOS_TYPE_ID) {
    // End-of-stream has no payload beyond type byte
  } else {
    store<u32>(ptr + 1, len);
    if (len > 0) {
      String.UTF8.encodeUnsafe(changetype<usize>(message), message.length, ptr + 1 + 4, false);
    }
  }

  return ptr;
}

/** Return a nil (type 0) pointer — used to signal success from set_param. */
export function nilPtr(): usize {
  const ptr = heap.alloc(1);
  store<i8>(ptr, NIL_TYPE_ID);
  return ptr;
}
