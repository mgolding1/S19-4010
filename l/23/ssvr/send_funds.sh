
# See: https://stackoverflow.com/questions/49413001/ethereum-eth-sendtransaction-via-php-curl

# From: 0xc0c4B94355fD676a29856008e625B51d1acD04eD
# To: Address: 0x1d217e902Bc1deB2e75D1Ec44bcAE03A1227a126
From=0xc0c4B94355fD676a29856008e625B51d1acD04eD
To=0x1d217e902Bc1deB2e75D1Ec44bcAE03A1227a126

# Orig: curl -H "Content-type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x35fa3c7edd23b23bd714fd075d243097e14ed937","to":"0xdab9a603ed3f1cf7b2b89f1cb1b57145e4828796","gas":"0x15f90","gasPrice":"0x430e23400","value":"0x9b6e64a8ec60000"}],"id":"1"}' http://127.0.0.1:8545
curl -H "Content-type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0xc0c4B94355fD676a29856008e625B51d1acD04eD","to":"0x1d217e902Bc1deB2e75D1Ec44bcAE03A1227a126","gas":"0x15f90","gasPrice":"0x430e23400","value":"0x9b6e64a8ec60000"}],"id":"1"}' http://127.0.0.1:9545

# Example Response: {"id":"1","jsonrpc":"2.0","result":"0xf1e96688831e4b7f1297dfc9d76c00c0ac950365c79cfa10c7690790fca145ba"}

