make
rm -rf data
mkdir data
mkdir data/node1
./build/bin/geth --datadir data/node1/ init genesis.json
./build/bin/geth --datadir data/node1 --ipcdisable --port 30301 --networkid 7 --http --http.addr 'localhost' --http.port 8545 --http.api 'admin,debug,eth,miner,net,personal,txpool,web3' --allow-insecure-unlock