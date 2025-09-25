# view-creator

./viewkit view init testdeploy

./viewkit view inspect testdeploy

./viewkit view add query "Log {address topics data transactionHash blockNumber}" --name testdeploy            

./viewkit view add sdl "type FilteredAndDecodedLogs @materialized(if: false) {transactionHash: String}" --name testdeploy

./viewkit view add lens --args '{"src":"address", "value":"0x1e3aA9fE4Ef01D3cB3189c129a49E3C03126C636"}' --label "filter" --url "https://raw.githubusercontent.com/shinzonetwork/wasm-bucket/main/bucket/filter_transaction/filter_transaction.wasm" --name testdeploy

./viewkit wallet generate  

./viewkit view deploy testdeploy --target devnet 
