
all: 

build:
	( cd main ; go build )

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
FR=../../Lectures/Lect-01

Transactions.html: Transactions.md
	markdown-cli --input=./Transactions.md --output=Transactions.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre Transactions.html ${FR}/css/hpost >/tmp/Transactions.html
	mv /tmp/Transactions.html ./Transactions.html

.PHONY: test
test:
	( cd addr ; go test )
	( cd block ; go test )
	( cd config ; go test )
	( cd hash ; go test )
	( cd lib ; go test )
	( cd merkle ; go test )
	( cd mine ; go test )
	( cd main ; make test )

