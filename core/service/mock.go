package service

const GQL = `
mutation {
  b1: create_Block(input: {
    hash: "0xabc1",
    number: 1001,
    timestamp: "2023-07-01T12:00:00Z",
    parentHash: "0x0000",
    difficulty: "1000000",
    gasUsed: "21000",
    gasLimit: "8000000",
    nonce: "0x01",
    miner: "0xMiner1",
    size: "1234",
    stateRoot: "0xState1",
    sha3Uncles: "0xUncle1",
    transactionsRoot: "0xTxRoot1",
    receiptsRoot: "0xReceipt1",
    extraData: "0xData1"
  }) { hash }

  t1: create_Transaction(input: {
    hash: "0xtx1",
    blockHash: "0xabc1",
    blockNumber: 1001,
    from: "0xAlice",
    to: "0xBob",
    value: "100",
    gasUsed: "21000",
    gasPrice: "1000000000",
    inputData: "0x",
    nonce: "1",
    transactionIndex: "0",
    r: "0xr",
    s: "0xs",
    v: "0x1b"
  }) { hash }

  l1: create_Log(input: {
    address: "0xContract1",
    topics: ["0xTopic1"],
    data: "0xdata1",
    transactionHash: "0xtx1",
    blockHash: "0xabc1",
    blockNumber: 1001,
    transactionIndex: "0",
    logIndex: "0",
    removed: "false"
  }) { address }

  e1: create_Event(input: {
    contractAddress: "0xContract1",
    eventName: "Transfer",
    parameters: "{\"from\":\"0xAlice\",\"to\":\"0xBob\",\"value\":\"100\"}",
    transactionHash: "0xtx1",
    blockHash: "0xabc1",
    blockNumber: 1001,
    transactionIndex: "0",
    logIndex: "0"
  }) { eventName }

  b2: create_Block(input: {
    hash: "0xabc2",
    number: 1002,
    timestamp: "2023-07-01T12:10:00Z",
    parentHash: "0xabc1",
    difficulty: "1001000",
    gasUsed: "31000",
    gasLimit: "8000000",
    nonce: "0x02",
    miner: "0xMiner2",
    size: "1240",
    stateRoot: "0xState2",
    sha3Uncles: "0xUncle2",
    transactionsRoot: "0xTxRoot2",
    receiptsRoot: "0xReceipt2",
    extraData: "0xData2"
  }) { hash }

  t2: create_Transaction(input: {
    hash: "0xtx2",
    blockHash: "0xabc2",
    blockNumber: 1002,
    from: "0xCarol",
    to: "0xDave",
    value: "500",
    gasUsed: "22000",
    gasPrice: "1100000000",
    inputData: "0xabc",
    nonce: "2",
    transactionIndex: "0",
    r: "0xr",
    s: "0xs",
    v: "0x1b"
  }) { hash }

  t3: create_Transaction(input: {
    hash: "0xtx3",
    blockHash: "0xabc2",
    blockNumber: 1002,
    from: "0xEve",
    to: "0xFrank",
    value: "1000",
    gasUsed: "24000",
    gasPrice: "1050000000",
    inputData: "0xdef",
    nonce: "3",
    transactionIndex: "1",
    r: "0xr",
    s: "0xs",
    v: "0x1c"
  }) { hash }

  l2: create_Log(input: {
    address: "0xContract2",
    topics: ["0xTopic2", "0xTopic3"],
    data: "0xdata2",
    transactionHash: "0xtx2",
    blockHash: "0xabc2",
    blockNumber: 1002,
    transactionIndex: "0",
    logIndex: "0",
    removed: "false"
  }) { address }

  e2: create_Event(input: {
    contractAddress: "0xContract2",
    eventName: "Approval",
    parameters: "{\"owner\":\"0xCarol\",\"spender\":\"0xDave\",\"value\":\"500\"}",
    transactionHash: "0xtx2",
    blockHash: "0xabc2",
    blockNumber: 1002,
    transactionIndex: "0",
    logIndex: "0"
  }) { eventName }

  b3: create_Block(input: {
    hash: "0xabc3",
    number: 1003,
    timestamp: "2023-07-01T12:20:00Z",
    parentHash: "0xabc2",
    difficulty: "1002000",
    gasUsed: "33000",
    gasLimit: "8000000",
    nonce: "0x03",
    miner: "0xMiner3",
    size: "1250",
    stateRoot: "0xState3",
    sha3Uncles: "0xUncle3",
    transactionsRoot: "0xTxRoot3",
    receiptsRoot: "0xReceipt3",
    extraData: "0xData3"
  }) { hash }

  t4: create_Transaction(input: {
    hash: "0xtx4",
    blockHash: "0xabc3",
    blockNumber: 1003,
    from: "0xGeorge",
    to: "0xHelen",
    value: "750",
    gasUsed: "25000",
    gasPrice: "950000000",
    inputData: "0xghi",
    nonce: "4",
    transactionIndex: "0",
    r: "0xr",
    s: "0xs",
    v: "0x1b"
  }) { hash }

  l3: create_Log(input: {
    address: "0xContract3",
    topics: ["0xTopic4"],
    data: "0xdata3",
    transactionHash: "0xtx4",
    blockHash: "0xabc3",
    blockNumber: 1003,
    transactionIndex: "0",
    logIndex: "0",
    removed: "false"
  }) { address }

  e3: create_Event(input: {
    contractAddress: "0xContract3",
    eventName: "Mint",
    parameters: "{\"minter\":\"0xGeorge\",\"value\":\"750\"}",
    transactionHash: "0xtx4",
    blockHash: "0xabc3",
    blockNumber: 1003,
    transactionIndex: "0",
    logIndex: "0"
  }) { eventName }
}
`
