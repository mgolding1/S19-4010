
# If you are working on Windows

You will need gcc.  I instaled from: http://tdm-gcc.tdragon.net/

# if you are working on Mac

You will need to have 'brew' installed - google 'install brew' and then install:

```
	brew install gmp openssl llvm
```

# Other notes
 
Also it appears that my original instructions for hw-5 are missing the --addr parameter that you have to use when
validating a signature. 

# Linking probblem

You have to delete a package from go-ethereum

In /Users/`Your User Name`/go/src/github.com/ethereum/go-ethereum/vendor/github.com and in
/Users/`Your User Name`/go/src/github.com/Univ-Wyo-Education/Blockchain-4010-Fall-2018/vendor/github.com

find the directory `pborman` and remove the `uuid` package.  That is the directory under `pborman`.

Then verify that you have this package in `~/go/src/github.com/pborman/uuid` -- to do that
create the directory `~/go/src/github.com/pborman`
cd to it, then 

```
go get https://github.com/pborman/uuid.git
```

