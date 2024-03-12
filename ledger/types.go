package ledger

import "crypto/ecdsa"

type Person struct {
	Name       string
	Balance    int
	PrivateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

type Transaction struct {
	Sender    *Person
	Receiver  *Person
	Amount    int
	Signature []byte
}

type Pot struct {
	Total int
}

type Ledger struct {
	Pot          Pot
	Transactions []Transaction
	Participants []*Person
}
