# Decoder Lens

Example usage with `viewkit` CLI:

```bash
# create a view with name "decoder"
./build/viewkit view init decoder

# basic query to pull in raw log data from Ethereum mainnet, with some nested transaction fields
./build/viewkit view add query "Ethereum__Mainnet__Log { address topics data blockNumber transaction { hash from to } }" --name decoder

# an SDL that returns transformed log data with event names from decoded ABI
./build/viewkit view add sdl "type DecodedLog @materialized(if: false) {                                                               
hash: String
blockNumber: Int
from: String
to: String
logAddress: String
event: String
signature: String
}" --name decoder

# attaching a recently-built lens with ERC-20 ABI
./build/viewkit view add lens \
--name decoder \
--label decoder \
--path ./as/build/decode_log/decode_log.wasm \
--args '{"abi": "[{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false}]},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false}]}]"}'

# deploy to a local DefraDB with mock data
./build/viewkit view deploy decoder --target local

# deploy to ShinzoHub, so that a host can pick it up and run it against devnet data
./build/viewkit view deploy decoder --target devnet --rpc http://rpc.devnet.shinzo.network:8545
```
