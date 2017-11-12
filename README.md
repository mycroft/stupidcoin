stupidcoin
==========

Specs
-----

### Blockchain

```go
type Blockchain struct {
    last_index uint64
    blocks     []*Block
}
```

### Blocks

```go
type Block struct {
    index     uint64
    last_hash []byte
    timestamp uint64
    hash      []byte
}
```

### Transactions

```go
type TxInput struct {
    txhash []byte
    script []byte
}

type TxOutput struct {
    script []byte
    amount float64
}

type Transaction struct {
    hash    []byte
    inputs  []TxInput
    outputs []TxOutput
}
```

Scripts
-------

### Inputs:

Contains a script, which will complete unlock scripts.

### Output:

Must contain unlock script.
An unlock script will contains operation, which, played with input script, will return true when all operations are played and stack is empty.

### Execution

Output script will be added to input script in order to unlock this output script, and allow creating new output scripts.

### Scripts samples

## P2PK: Pay to Publickey

Input:

- signature

Output

- pubkey OP_CHECKSIG

## P2PKH: Pay to PublickeyHash

Input:

- signature publickey

Output

- OP_DUP OP_HASH160 pubkeyhash OP_EQUALVERIFY OP_CHECKSIG


Api
---

GET /mine

POST /txn/add