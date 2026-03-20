export const ERC20_ABI = JSON.stringify([
  {
    type: "event",
    name: "Transfer",
    inputs: [
      { name: "from", type: "address", indexed: true },
      { name: "to", type: "address", indexed: true },
      { name: "value", type: "uint256", indexed: false },
    ],
  },
  {
    type: "event",
    name: "Approval",
    inputs: [
      { name: "owner", type: "address", indexed: true },
      { name: "spender", type: "address", indexed: true },
      { name: "value", type: "uint256", indexed: false },
    ],
  },
]);

// keccak256("Transfer(address,address,uint256)")
const TRANSFER_TOPIC0 = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef";

export const SAMPLE_LOG = JSON.stringify({
  address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
  topics: [
    TRANSFER_TOPIC0,
    "0x000000000000000000000000ab5801a7d398351b8be11c439e05c5b3259aec9b",
    "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
  ],
  data: "0x0000000000000000000000000000000000000000000000000000000077359400",
  blockNumber: 18500000,
  transaction: {
    hash: "0xabc123",
    from: "0xsender",
    to: "0xreceiver",
  },
});
