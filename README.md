# view-creator

go run ./cmd/viewkit/ view deploy testdeploy --target local 


go build -o viewkit ./cmd/viewkit


curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionByHash",
    "params": ["0x4bec7f625df55ba1e11e3715df8dd8ca72f2fd43e9f74cd8a3497c802ccb4537"],
    "id": 1
  }'

curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"eth_getTransactionReceipt",
    "params":["0x4bec7f625df55ba1e11e3715df8dd8ca72f2fd43e9f74cd8a3497c802ccb4537"],
    "id":1
  }'

  view init testdeploy
  view init inspect