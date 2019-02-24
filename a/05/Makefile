
all: Wallet.html build

build:
	( cd bsvr/main ; go build )

# You can get markdown-cli and compile it yourself.
#
# $ cd ~/go/src/github.com/
# $ mkdir pschlump
# $ cd pschlump
# $ git pull https://github.com/pschlump/markdown-cli.git
# $ go get
# $ cd markdown-cli
# $ go build
#
FR=../../l/01

Wallet.html: Wallet.md
	markdown-cli --input=./Wallet.md --output=Wallet.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre Wallet.html ${FR}/css/hpost >/tmp/Wallet.html
	mv /tmp/Wallet.html ./Wallet.html

.PHONY: test
test:
	( cd bsvr/addr ; go test )
	( cd bsvr/block ; go test )
	( cd bsvr/config ; go test )
	( cd bsvr/hash ; go test )
	( cd bsvr/lib ; go test )
	( cd bsvr/merkle ; go test )
	( cd bsvr/mine ; go test )
	( cd bsvr/main ; make test )

.PHONY: run_bsvr
run_bsvr:
	( cd bsvr/main ; ./main )


	
