Pt 1: Merkle Trees : Assignment 3
==============================

One of the major security and validity checks that blockchains
do is using Merkle trees.

In this assignment you will implement a Merkle tree hash.

When you pull code from git you should have a `./A-03/hash` directory.
This is a copy of Assignment 2's `./A-02/hash`.  The new code that you will
be working on is in `./merkle`.


### Pseudo Code

1. Create a slice to hold the hashes of the leaves.  Each leaf hash is a `[]byte`.  So make the data type `[][]byte`.
	Make this slice of slice of byte then length of the data.  That would be `len(data)`.  Let's call this `hTmp`.
2. For each data block
	1. Calculate a hash for the data block using `hash.HashOf()`.
	2. Save this in the slice created in (1) above.
3. Create a `[][]byte` slice to hold the intermediate hashes in the tree.
	This will need to be no more than `len(data)/2+1` in length.  The plus 1 is so that 0 blocks of hasing or an odd
	number of blocks will have enough space.  Let's call this `hMid`.
4. Declare a variable `ln`, and set it to  `len(data)/2+1`
5. While `ln >= 1` (Hint: the language only has `for` loops with lots of different ways of doing it)
	1. For each pair of hashes (if you have an odd number just use the single hash)
		- Calculate the hash of the pair using `hash.Keccak256()`.  It takes a variable number of arguments so you can
		  pass 1 or 2 arguments to it.
		- Append this to `hMid`.
	2. Replace hTmp with hMid
	3. Recalculate `ln` set it to `len(hTmp)/2`
	4. Generate a new empty hMid of allocated space of `len(hTmp)/2`.
6. Return `hTmp[0]`


# Submit

1. Your code, ./merkle/merkle.go.
2. Any additional test cases that you created.
3. Your prove that this works.

# References

1. [Wikipedia has a nice discussion](https://en.wikipedia.org/wiki/Merkle_tree)
2. [Another explanation of Merkle Trees - with more details](https://brilliant.org/wiki/merkle-tree/)

Pt 2: Basic Client Server
==============================

This part of the homework is to add a very simple API call to a client server.  The server is in ./simple-server, the sample client is in
./simple-client.

Modify the sever to have a `/double` that will take a value that it is passed on a `GET` call and double that value.  The return the
doubled value to the client in a JSON format.

So...

```
	http://localhost:3000/double?value=12
```

will return

```
	{"double":24}
```

You will need to set the MIME type of the return to JSON with

```
	xyzzy
```

Use curl or wget to test this call.

```
	wget 'http://localhost:3000/double?value=100'
```

or

```
	curl 'http://localhost:3000/double?value=100'
```

