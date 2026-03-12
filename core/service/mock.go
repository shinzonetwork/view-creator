package service

import (
	"fmt"
	"strings"
)

const GQL_BLOCKS = `mutation {
  b1: add_Ethereum__Mainnet__Block(input: {
    hash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    number: 18500000,
    timestamp: "2023-11-01T12:00:00Z",
    parentHash: "0x3c2a2d4a5c1e8f6b9d7e0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d",
    difficulty: "0",
    gasUsed: "12000000",
    gasLimit: "30000000",
    nonce: "0x0000000000000000",
    miner: "0x388C818CA8B9251b393131C08a736A67ccB19297",
    size: "1200",
    stateRoot: "0xabc123def456789abc123def456789abc123def456789abc123def456789abcd",
    sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    transactionsRoot: "0x5678abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456",
    receiptsRoot: "0x9876fedcba0987654321fedcba0987654321fedcba0987654321fedcba098765",
    extraData: "0x"
  }) { hash }

  b2: add_Ethereum__Mainnet__Block(input: {
    hash: "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
    number: 18500001,
    timestamp: "2023-11-01T12:00:12Z",
    parentHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    difficulty: "0",
    gasUsed: "15000000",
    gasLimit: "30000000",
    nonce: "0x0000000000000000",
    miner: "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326",
    size: "1350",
    stateRoot: "0xdef456789abc123def456789abc123def456789abc123def456789abc123def4",
    sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    transactionsRoot: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
    receiptsRoot: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    extraData: "0x"
  }) { hash }

  b3: add_Ethereum__Mainnet__Block(input: {
    hash: "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
    number: 18500002,
    timestamp: "2023-11-01T12:00:24Z",
    parentHash: "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
    difficulty: "0",
    gasUsed: "18000000",
    gasLimit: "30000000",
    nonce: "0x0000000000000000",
    miner: "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
    size: "1500",
    stateRoot: "0x789abc123def456789abc123def456789abc123def456789abc123def456789a",
    sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    transactionsRoot: "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
    receiptsRoot: "0x0987654321fedcba0987654321fedcba0987654321fedcba0987654321fedcba",
    extraData: "0x"
  }) { hash }
}`

const GQL_TRANSACTIONS = `mutation {
  t1_1: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xaa1c8b56e2f34d79e0b1a2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3",
    blockHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    blockNumber: 18500000,
    from: "0xDFd5293D8e347dFe59E90eFd55b2956a1343963d",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "65000", gasPrice: "30000000000",
    input: "0xa9059cbb00000000000000000000000028C6c06298d514Db089934071355E5743bf21d6000000000000000000000000000000000000000000000000000000000003d0900",
    nonce: "100", transactionIndex: 0,
    r: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    s: "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321", v: "0x1c"
  }) { hash }

  t1_2: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xbb2d9c67f3e45a80f1c2b3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4",
    blockHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    blockNumber: 18500000,
    from: "0x28C6c06298d514Db089934071355E5743bf21d60",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "52000", gasPrice: "28000000000",
    input: "0xa9059cbb000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA96045000000000000000000000000000000000000000000000000000000174876e800",
    nonce: "5012", transactionIndex: 1,
    r: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
    s: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", v: "0x1b"
  }) { hash }

  t1_3: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xcc3eab78a4f56b91a2d3c4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5",
    blockHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    blockNumber: 18500000,
    from: "0x47ac0Fb4F2D84898e4D9E7b4DaB3C24507a6D503",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "48000", gasPrice: "32000000000",
    input: "0x095ea7b3000000000000000000000000E592427A0AEce92De3Edee1F18E0157C05861564ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
    nonce: "782", transactionIndex: 2,
    r: "0x567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234",
    s: "0xcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab", v: "0x1c"
  }) { hash }

  t1_4: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xdd4fbc89b5a67ca2b3e4d5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6",
    blockHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    blockNumber: 18500000,
    from: "0x56Eddb7aa87536c09CCc2793473599fD21A8b17F",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "65000", gasPrice: "35000000000",
    input: "0xa9059cbb000000000000000000000000Dac17F958D2ee523a2206206994597C13D831ec70000000000000000000000000000000000000000000000000000000005f5e100",
    nonce: "341", transactionIndex: 3,
    r: "0x890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678",
    s: "0xef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", v: "0x1b"
  }) { hash }

  t1_5: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xee5acd9ac6b78db3c4f5e6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7",
    blockHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    blockNumber: 18500000,
    from: "0x21a31Ee1afC51d94C2eFcCAa2092aD1028285549",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "65000", gasPrice: "29000000000",
    input: "0xa9059cbb000000000000000000000000BE0eB53F46cd790Cd13851d5EFf43D12404d33E80000000000000000000000000000000000000000000000000000001caab5c3b3",
    nonce: "456", transactionIndex: 4,
    r: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
    s: "0x567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234", v: "0x1c"
  }) { hash }

  t1_6: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xab01usdt1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d",
    blockHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    blockNumber: 18500000,
    from: "0xDFd5293D8e347dFe59E90eFd55b2956a1343963d",
    to: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
    value: "0", gasUsed: "63000", gasPrice: "30000000000",
    input: "0xa9059cbb00000000000000000000000028C6c06298d514Db089934071355E5743bf21d600000000000000000000000000000000000000000000000000000000005f5e100",
    nonce: "101", transactionIndex: 5,
    r: "0xaaaa567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    s: "0xbbbbba0987654321fedcba0987654321fedcba0987654321fedcba0987654321", v: "0x1c"
  }) { hash }

  t1_7: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xcd02weth2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e",
    blockHash: "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
    blockNumber: 18500000,
    from: "0x28C6c06298d514Db089934071355E5743bf21d60",
    to: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    value: "1000000000000000000", gasUsed: "45000", gasPrice: "30000000000",
    input: "0xd0e30db0",
    nonce: "5013", transactionIndex: 6,
    r: "0xcccc567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    s: "0xddddba0987654321fedcba0987654321fedcba0987654321fedcba0987654321", v: "0x1b"
  }) { hash }

  t2_1: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xff6bde0ab7c89ec4d5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8",
    blockHash: "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
    blockNumber: 18500001,
    from: "0xBE0eB53F46cd790Cd13851d5EFf43D12404d33E8",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "65000", gasPrice: "31000000000",
    input: "0xa9059cbb0000000000000000000000005041ed759Dd4aFc3a72b8192C143F72f4724081A00000000000000000000000000000000000000000000000000000002540be400",
    nonce: "9001", transactionIndex: 0,
    r: "0xdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd",
    s: "0x90abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678", v: "0x1c"
  }) { hash }

  t2_2: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0x007cef1bc8d9afd5e6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9",
    blockHash: "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
    blockNumber: 18500001,
    from: "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "52000", gasPrice: "33000000000",
    input: "0xa9059cbb000000000000000000000000F977814e90dA44bFA03b6295A0616a897441aceC000000000000000000000000000000000000000000000000000000003b9aca00",
    nonce: "221", transactionIndex: 1,
    r: "0x4567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234",
    s: "0xbcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab", v: "0x1b"
  }) { hash }

  t2_3: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0x118def2cd9eabae6f7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa",
    blockHash: "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
    blockNumber: 18500001,
    from: "0xF977814e90dA44bFA03b6295A0616a897441aceC",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "48000", gasPrice: "30000000000",
    input: "0x095ea7b300000000000000000000000068b3465833fb72A70ecDF485E0e4C7bD8665Fc4500000000000000000000000000000000000000000000000000000000ee6b2800",
    nonce: "1024", transactionIndex: 2,
    r: "0x7890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456",
    s: "0x34567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12", v: "0x1c"
  }) { hash }

  t2_4: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xef03dai03c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f",
    blockHash: "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
    blockNumber: 18500001,
    from: "0xBE0eB53F46cd790Cd13851d5EFf43D12404d33E8",
    to: "0x6B175474E89094C44Da98b954EedeAC495271d0F",
    value: "0", gasUsed: "60000", gasPrice: "31000000000",
    input: "0xa9059cbb000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA960450000000000000000000000000000000000000000000000056bc75e2d63100000",
    nonce: "9002", transactionIndex: 3,
    r: "0xeeee1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd",
    s: "0xffff0abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678", v: "0x1c"
  }) { hash }

  t2_5: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xa104weth4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a",
    blockHash: "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
    blockNumber: 18500001,
    from: "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
    to: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    value: "0", gasUsed: "55000", gasPrice: "32000000000",
    input: "0xa9059cbb000000000000000000000000F977814e90dA44bFA03b6295A0616a897441aceC00000000000000000000000000000000000000000000000029a2241af62c0000",
    nonce: "222", transactionIndex: 4,
    r: "0x11111234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd",
    s: "0x22220abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678", v: "0x1b"
  }) { hash }

  t3_1: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0x229ef03dea0bcbf7a8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa0b",
    blockHash: "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
    blockNumber: 18500002,
    from: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "65000", gasPrice: "27000000000",
    input: "0xa9059cbb000000000000000000000000DFd5293D8e347dFe59E90eFd55b2956a1343963d0000000000000000000000000000000000000000000000000000000000989680",
    nonce: "512", transactionIndex: 0,
    r: "0xcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
    s: "0x4567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234", v: "0x1b"
  }) { hash }

  t3_2: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0x33af014efb1cdca8b9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa0b1c",
    blockHash: "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
    blockNumber: 18500002,
    from: "0x5041ed759Dd4aFc3a72b8192C143F72f4724081A",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "65000", gasPrice: "28000000000",
    input: "0xa9059cbb00000000000000000000000047ac0Fb4F2D84898e4D9E7b4DaB3C24507a6D503000000000000000000000000000000000000000000000000000000012a05f200",
    nonce: "89", transactionIndex: 1,
    r: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    s: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", v: "0x1c"
  }) { hash }

  t3_3: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0x44b0125fac2eddb9c0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa0b1c2d",
    blockHash: "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
    blockNumber: 18500002,
    from: "0xE592427A0AEce92De3Edee1F18E0157C05861564",
    to: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    value: "0", gasUsed: "65000", gasPrice: "26000000000",
    input: "0xa9059cbb00000000000000000000000056Eddb7aa87536c09CCc2793473599fD21A8b17F0000000000000000000000000000000000000000000000000000000077359400",
    nonce: "3001", transactionIndex: 2,
    r: "0x890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678",
    s: "0xef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", v: "0x1b"
  }) { hash }

  t3_4: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xb205usdt5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b",
    blockHash: "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
    blockNumber: 18500002,
    from: "0x47ac0Fb4F2D84898e4D9E7b4DaB3C24507a6D503",
    to: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
    value: "0", gasUsed: "46000", gasPrice: "27000000000",
    input: "0x095ea7b3000000000000000000000000E592427A0AEce92De3Edee1F18E0157C05861564ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
    nonce: "783", transactionIndex: 3,
    r: "0x33331234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd",
    s: "0x44440abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678", v: "0x1c"
  }) { hash }

  t3_5: add_Ethereum__Mainnet__Transaction(input: {
    hash: "0xc306dai06f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c",
    blockHash: "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
    blockNumber: 18500002,
    from: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
    to: "0x6B175474E89094C44Da98b954EedeAC495271d0F",
    value: "0", gasUsed: "46000", gasPrice: "26000000000",
    input: "0x095ea7b300000000000000000000000068b3465833fb72A70ecDF485E0e4C7bD8665Fc45ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
    nonce: "513", transactionIndex: 4,
    r: "0x55551234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd",
    s: "0x66660abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678", v: "0x1b"
  }) { hash }
}`

type mockLog struct {
	alias            string
	address          string
	topics           []string
	data             string
	transactionHash  string
	blockHash        string
	blockNumber      int
	transactionIndex int
	logIndex         int
}

var mockLogs = []mockLog{
	{
		alias:   "l1_1_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x000000000000000000000000DFd5293D8e347dFe59E90eFd55b2956a1343963d", "0x00000000000000000000000028C6c06298d514Db089934071355E5743bf21d60"},
		data:             "0x00000000000000000000000000000000000000000000000000000000003d0900",
		transactionHash:  "0xaa1c8b56e2f34d79e0b1a2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3",
		blockHash:        "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
		blockNumber:      18500000,
		transactionIndex: 0,
		logIndex:         0,
	},
	{
		alias:   "l1_2_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x00000000000000000000000028C6c06298d514Db089934071355E5743bf21d60", "0x000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA96045"},
		data:             "0x000000000000000000000000000000000000000000000000000000174876e800",
		transactionHash:  "0xbb2d9c67f3e45a80f1c2b3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4",
		blockHash:        "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
		blockNumber:      18500000,
		transactionIndex: 1,
		logIndex:         1,
	},
	{
		alias:   "l1_3_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", "0x00000000000000000000000047ac0Fb4F2D84898e4D9E7b4DaB3C24507a6D503", "0x000000000000000000000000E592427A0AEce92De3Edee1F18E0157C05861564"},
		data:             "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		transactionHash:  "0xcc3eab78a4f56b91a2d3c4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5",
		blockHash:        "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
		blockNumber:      18500000,
		transactionIndex: 2,
		logIndex:         2,
	},
	{
		alias:   "l1_4_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x00000000000000000000000056Eddb7aa87536c09CCc2793473599fD21A8b17F", "0x000000000000000000000000Dac17F958D2ee523a2206206994597C13D831ec7"},
		data:             "0x0000000000000000000000000000000000000000000000000000000005f5e100",
		transactionHash:  "0xdd4fbc89b5a67ca2b3e4d5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6",
		blockHash:        "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
		blockNumber:      18500000,
		transactionIndex: 3,
		logIndex:         3,
	},
	{
		alias:   "l1_5_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x00000000000000000000000021a31Ee1afC51d94C2eFcCAa2092aD1028285549", "0x000000000000000000000000BE0eB53F46cd790Cd13851d5EFf43D12404d33E8"},
		data:             "0x0000000000000000000000000000000000000000000000000000001caab5c3b3",
		transactionHash:  "0xee5acd9ac6b78db3c4f5e6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7",
		blockHash:        "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
		blockNumber:      18500000,
		transactionIndex: 4,
		logIndex:         4,
	},
	{
		alias:   "l1_6_1",
		address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x000000000000000000000000DFd5293D8e347dFe59E90eFd55b2956a1343963d", "0x00000000000000000000000028C6c06298d514Db089934071355E5743bf21d60"},
		data:             "0x0000000000000000000000000000000000000000000000000000000005f5e100",
		transactionHash:  "0xab01usdt1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d",
		blockHash:        "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
		blockNumber:      18500000,
		transactionIndex: 5,
		logIndex:         5,
	},
	{
		alias:   "l1_7_1",
		address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		topics:  []string{"0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c", "0x00000000000000000000000028C6c06298d514Db089934071355E5743bf21d60"},
		data:             "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000",
		transactionHash:  "0xcd02weth2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e",
		blockHash:        "0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd",
		blockNumber:      18500000,
		transactionIndex: 6,
		logIndex:         6,
	},
	{
		alias:   "l2_1_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x000000000000000000000000BE0eB53F46cd790Cd13851d5EFf43D12404d33E8", "0x0000000000000000000000005041ed759Dd4aFc3a72b8192C143F72f4724081A"},
		data:             "0x00000000000000000000000000000000000000000000000000000002540be400",
		transactionHash:  "0xff6bde0ab7c89ec4d5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8",
		blockHash:        "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
		blockNumber:      18500001,
		transactionIndex: 0,
		logIndex:         0,
	},
	{
		alias:   "l2_2_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x0000000000000000000000003fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD", "0x000000000000000000000000F977814e90dA44bFA03b6295A0616a897441aceC"},
		data:             "0x000000000000000000000000000000000000000000000000000000003b9aca00",
		transactionHash:  "0x007cef1bc8d9afd5e6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9",
		blockHash:        "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
		blockNumber:      18500001,
		transactionIndex: 1,
		logIndex:         1,
	},
	{
		alias:   "l2_3_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", "0x000000000000000000000000F977814e90dA44bFA03b6295A0616a897441aceC", "0x00000000000000000000000068b3465833fb72A70ecDF485E0e4C7bD8665Fc45"},
		data:             "0x00000000000000000000000000000000000000000000000000000000ee6b2800",
		transactionHash:  "0x118def2cd9eabae6f7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa",
		blockHash:        "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
		blockNumber:      18500001,
		transactionIndex: 2,
		logIndex:         2,
	},
	{
		alias:   "l2_4_1",
		address: "0x6B175474E89094C44Da98b954EedeAC495271d0F",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x000000000000000000000000BE0eB53F46cd790Cd13851d5EFf43D12404d33E8", "0x000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA96045"},
		data:             "0x0000000000000000000000000000000000000000000000056bc75e2d63100000",
		transactionHash:  "0xef03dai03c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f",
		blockHash:        "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
		blockNumber:      18500001,
		transactionIndex: 3,
		logIndex:         3,
	},
	{
		alias:   "l2_5_1",
		address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x0000000000000000000000003fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD", "0x000000000000000000000000F977814e90dA44bFA03b6295A0616a897441aceC"},
		data:             "0x00000000000000000000000000000000000000000000000029a2241af62c0000",
		transactionHash:  "0xa104weth4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a",
		blockHash:        "0x5f4b4865521288f7a48ef2f95bba69ea239f9d2b3369d6f96eb9f2de826b2cee",
		blockNumber:      18500001,
		transactionIndex: 4,
		logIndex:         4,
	},
	{
		alias:   "l3_1_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA96045", "0x000000000000000000000000DFd5293D8e347dFe59E90eFd55b2956a1343963d"},
		data:             "0x0000000000000000000000000000000000000000000000000000000000989680",
		transactionHash:  "0x229ef03dea0bcbf7a8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa0b",
		blockHash:        "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
		blockNumber:      18500002,
		transactionIndex: 0,
		logIndex:         0,
	},
	{
		alias:   "l3_2_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x0000000000000000000000005041ed759Dd4aFc3a72b8192C143F72f4724081A", "0x00000000000000000000000047ac0Fb4F2D84898e4D9E7b4DaB3C24507a6D503"},
		data:             "0x000000000000000000000000000000000000000000000000000000012a05f200",
		transactionHash:  "0x33af014efb1cdca8b9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa0b1c",
		blockHash:        "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
		blockNumber:      18500002,
		transactionIndex: 1,
		logIndex:         1,
	},
	{
		alias:   "l3_3_1",
		address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x000000000000000000000000E592427A0AEce92De3Edee1F18E0157C05861564", "0x00000000000000000000000056Eddb7aa87536c09CCc2793473599fD21A8b17F"},
		data:             "0x0000000000000000000000000000000000000000000000000000000077359400",
		transactionHash:  "0x44b0125fac2eddb9c0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9fa0b1c2d",
		blockHash:        "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
		blockNumber:      18500002,
		transactionIndex: 2,
		logIndex:         2,
	},
	{
		alias:   "l3_4_1",
		address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		topics:  []string{"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", "0x00000000000000000000000047ac0Fb4F2D84898e4D9E7b4DaB3C24507a6D503", "0x000000000000000000000000E592427A0AEce92De3Edee1F18E0157C05861564"},
		data:             "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		transactionHash:  "0xb205usdt5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b",
		blockHash:        "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
		blockNumber:      18500002,
		transactionIndex: 3,
		logIndex:         3,
	},
	{
		alias:   "l3_5_1",
		address: "0x6B175474E89094C44Da98b954EedeAC495271d0F",
		topics:  []string{"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", "0x000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA96045", "0x00000000000000000000000068b3465833fb72A70ecDF485E0e4C7bD8665Fc45"},
		data:             "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		transactionHash:  "0xc306dai06f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c",
		blockHash:        "0x6a5c5976632399a8b59f3fa6ccb7aefb34af0e3c4480e7fa7fc0a3ef937c3dff",
		blockNumber:      18500002,
		transactionIndex: 4,
		logIndex:         4,
	},
}

func BuildLogsMutation(blockDocIDs, txDocIDs map[string]string) string {
	var sb strings.Builder
	sb.WriteString("mutation {\n")
	for _, log := range mockLogs {
		quotedTopics := make([]string, len(log.topics))
		for i, t := range log.topics {
			quotedTopics[i] = fmt.Sprintf("%q", t)
		}
		topicsStr := "[" + strings.Join(quotedTopics, ", ") + "]"

		txDocID := txDocIDs[log.transactionHash]
		blockDocID := blockDocIDs[log.blockHash]

		sb.WriteString(fmt.Sprintf(
			"  %s: add_Ethereum__Mainnet__Log(input: {\n"+
				"    address: %q,\n"+
				"    topics: %s,\n"+
				"    data: %q,\n"+
				"    transactionHash: %q,\n"+
				"    blockHash: %q,\n"+
				"    blockNumber: %d,\n"+
				"    transactionIndex: %d,\n"+
				"    logIndex: %d,\n"+
				"    removed: \"false\",\n"+
				"    transaction: %q,\n"+
				"    block: %q\n"+
				"  }) { address }\n\n",
			log.alias, log.address, topicsStr, log.data,
			log.transactionHash, log.blockHash,
			log.blockNumber, log.transactionIndex, log.logIndex,
			txDocID, blockDocID,
		))
	}
	sb.WriteString("}")
	return sb.String()
}
