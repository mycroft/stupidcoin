stupidcoin
==========

TODO
----

- Créer une paire clef publique/privée, la sauvegarder dans le wallet.
    - Support de plusieurs clefs
    - Dump
- Blocks
- Transactions
- Scripts


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
