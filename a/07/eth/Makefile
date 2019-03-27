
all: \
	./contracts/SignedData_sol_SignedData.abi \
	./contracts/SignedDataVersion01_sol_SignedDataVersion01.abi 

truffle_test:
	truffle test --network development

test3_migrate:
	./unlock-6ff-3.sh
	truffle migrate --reset --network test3 | tee m.out
	upd-gcall <m.out
	save-old.sh

lint:
	solium --file contracts/SignedData.sol
	solium --file contracts/SignedDataVersion01.sol



./contracts/SignedData_sol_SignedData.abi: ./contracts/SignedData.sol
	( cd contracts ; solc openzeppelin-solidity/=../openzeppelin-solidity/ --allow-paths . --overwrite --abi -o . SignedData.sol ; mv SignedData.abi SignedData_sol_SignedData.abi )
	( cd contracts ; /Users/corwin/bin/check-json-syntax -p <SignedData_sol_SignedData.abi >,tmp ; mv ,tmp SignedData_sol_SignedData.abi )
	( cd contracts ; solc openzeppelin-solidity/=../openzeppelin-solidity/ --allow-paths . --overwrite --bin -o . SignedData.sol ; mv SignedData.bin SignedData_sol_SignedData.bin )
	( cd contracts ; solc openzeppelin-solidity/=../openzeppelin-solidity/ --allow-paths . --overwrite --bin-runtime -o . SignedData.sol )
	( cd contracts ; abigen --abi SignedData_sol_SignedData.abi --pkg SignedData --out SignedData.go )
	mkdir -p ./lib/SignedData
	cp ./contracts/SignedData.go ./lib/SignedData
	cp ./contracts/SignedData_sol_SignedData.abi ./abi
	( cd ./lib/SignedData ; go build )


./contracts/SignedDataVersion01_sol_SignedDataVersion01.abi: ./contracts/SignedDataVersion01.sol
	( cd contracts ; solc openzeppelin-solidity/=../openzeppelin-solidity/ --allow-paths . --overwrite --abi -o . SignedDataVersion01.sol ; mv SignedDataVersion01.abi SignedDataVersion01_sol_SignedDataVersion01.abi )
	( cd contracts ; /Users/corwin/bin/check-json-syntax -p <SignedDataVersion01_sol_SignedDataVersion01.abi >,tmp ; mv ,tmp SignedDataVersion01_sol_SignedDataVersion01.abi )
	( cd contracts ; solc openzeppelin-solidity/=../openzeppelin-solidity/ --allow-paths . --overwrite --bin -o . SignedDataVersion01.sol ; mv SignedDataVersion01.bin SignedDataVersion01_sol_SignedDataVersion01.bin )
	( cd contracts ; solc openzeppelin-solidity/=../openzeppelin-solidity/ --allow-paths . --overwrite --bin-runtime -o . SignedDataVersion01.sol )
	( cd contracts ; abigen --abi SignedDataVersion01_sol_SignedDataVersion01.abi --pkg SignedDataVersion01 --out SignedDataVersion01.go )
	mkdir -p ./lib/SignedDataVersion01
	cp ./contracts/SignedDataVersion01.go ./lib/SignedDataVersion01
	cp ./contracts/SignedDataVersion01_sol_SignedDataVersion01.abi ./abi
	( cd ./lib/SignedDataVersion01 ; go build )







migrate:
	truffle migrate --reset	

philip_migrate:
	date
	/Users/corwin/Projects/eth158/unlock-c2a.sh
	truffle migrate --reset --network local

upd_addr:
	upd-gcall <m.out
	save-old.sh

philip_test:
	date
	/Users/corwin/Projects/eth158/unlock-c2a.sh
	truffle test --network local

philip_unlock:
	/Users/corwin/Projects/eth158/unlock-c2a.sh

t_dev:
	truffle develop log

test:
	truffle test
.PHONY: test

test_local:
	truffle test --network local
.PHONY: test_local

