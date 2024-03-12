package ledger

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
)

func (p Person) String() string {
	return fmt.Sprintf("Name: %s, Balance: %d", p.Name, p.Balance)
}

func (p *Person) GenerateKeyPair() error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	p.PrivateKey = *privateKey
	p.PublicKey = privateKey.PublicKey
	return nil
}

func (t *Transaction) Sign() error {
	data := fmt.Sprintf("%s%s%d", t.Sender.Name, t.Receiver.Name, t.Amount)
	hashedData := sha256.Sum256([]byte(data))

	sig, err := ecdsa.SignASN1(rand.Reader, &t.Sender.PrivateKey, hashedData[:])
	if err != nil {
		return err
	}

	t.Signature = sig
	return nil
}

func (t Transaction) Verify() bool {
	data := fmt.Sprintf("%s%s%d", t.Sender.Name, t.Receiver.Name, t.Amount)
	hashedData := sha256.Sum256([]byte(data))

	return ecdsa.VerifyASN1(&t.Sender.PublicKey, hashedData[:], t.Signature)
}

func (t Transaction) String() string {
	return fmt.Sprintf("Sender: %s, Receiver: %s, Amount: %d", t.Sender.Name, t.Receiver.Name, t.Amount)
}

func (p Pot) String() string {
	return fmt.Sprintf("%d", p.Total)
}

func NewLedger() *Ledger {
	ledger := &Ledger{
		Pot:          Pot{Total: 100},
		Transactions: []Transaction{},
		Participants: []*Person{
			{Name: "Felipe", Balance: 0},
			{Name: "John", Balance: 0},
			{Name: "Satoshi", Balance: 0},
			{Name: "Jane", Balance: 0},
		},
	}

	fmt.Println("Generating ECDSA key pair for each participant")
	for _, participant := range ledger.Participants {
		err := participant.GenerateKeyPair()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Participant: %s\nPublic Key: X: %x Y: %x\n\n", participant.Name, participant.PublicKey.X.Bytes(), participant.PublicKey.Y.Bytes())
	}

	return ledger
}

func (l *Ledger) RunInitialDistribution() {
	distAmount := 20
	for _, p := range l.Participants {
		l.Pot.Total -= distAmount
		p.Balance += distAmount
	}
}

func (l *Ledger) addTransaction(sender *Person, receiver *Person, amount int) error {
	newTransaction := Transaction{Sender: sender, Receiver: receiver, Amount: amount}
	err := newTransaction.Sign()
	if err != nil {
		return err
	}

	l.Transactions = append(l.Transactions, newTransaction)
	return nil
}

func (l *Ledger) updateBalancePostTransaction(sender *Person, receiver *Person, amount int) {
	sender.Balance -= amount
	receiver.Balance += amount
}

func validateTransaction(sender *Person, receiver *Person, amount int) error {
	if sender.Balance < amount {
		return errors.New("sender has insufficient funds")
	}

	if amount < 0 {
		return errors.New("cannot send negative amounts")
	}

	if sender == receiver {
		return errors.New("cannot send to self")
	}

	return nil
}

func (l *Ledger) GetParticipantByName(name string) (*Person, error) {
	for _, p := range l.Participants {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, errors.New("Person does not belong to the ledger")
}

func (l *Ledger) SimulateTransaction(sender *Person, receiver *Person, amount int) error {
	err := validateTransaction(sender, receiver, amount)
	if err != nil {
		return err
	}

	newTransaction := &Transaction{Sender: sender, Receiver: receiver, Amount: amount}
	err = newTransaction.Sign()
	if err != nil {
		return err
	}

	isValidTransaction := newTransaction.Verify()
	if !isValidTransaction {
		return errors.New("failed to verify transaction")
	}

	l.updateBalancePostTransaction(sender, receiver, amount)
	l.Transactions = append(l.Transactions, *newTransaction)
	return nil
}

func (l *Ledger) Print() {
	fmt.Printf("Ledger:\n\tPot: %v\n", l.Pot)
	fmt.Println("\tParticipants:")
	for _, p := range l.Participants {
		fmt.Printf("\t\t%v\n", p)
	}
	fmt.Println("\tTransactions:")
	if len(l.Transactions) == 0 {
		fmt.Println("\t\t0")
	} else {

		for idx, t := range l.Transactions {
			fmt.Printf("\t\ttx-%d: %v\n", idx, t)
		}

	}
}

func RunSimulation() {
	fmt.Print("Initializing new communal ledger\n\n")
	communalLedger := NewLedger()
	communalLedger.Print()

	fmt.Println("Distributing 20 coins to each participant")
	communalLedger.RunInitialDistribution()
	communalLedger.Print()
	fmt.Println("Simulating transactions")
	var err error

	felipe, err := communalLedger.GetParticipantByName("Felipe")
	if err != nil {
		fmt.Println(err)
	}

	satoshi, err := communalLedger.GetParticipantByName("Satoshi")
	if err != nil {
		fmt.Println(err)
	}

	if felipe != nil && satoshi != nil {
		fmt.Println("Felipe is attempting to send Satoshi 5 coins")
		err = communalLedger.SimulateTransaction(felipe, satoshi, 5)
		if err != nil {
			fmt.Println("Transaction error:", err)
		} else {
			tx0 := communalLedger.Transactions[0]
			fmt.Println("Transaction successful\nDetails:")
			fmt.Printf("\tSender: %v\n\tReceiver: %v\n\tAmount: %d\n\tSignature: %x\n", tx0.Sender.Name, tx0.Receiver.Name, tx0.Amount, tx0.Signature)
			communalLedger.Print()
		}
	}

	jane, err := communalLedger.GetParticipantByName("Jane")
	if err != nil {
		fmt.Println(err)
	}

	if jane != nil && satoshi != nil {
		fmt.Println("Satoshi is attempting to send Jane 15 coins")
		err = communalLedger.SimulateTransaction(satoshi, jane, 15)
		if err != nil {
			fmt.Println("Transaction error:", err)
		} else {
			tx1 := communalLedger.Transactions[1]
			fmt.Println("Transaction successful\nDetails:")
			fmt.Printf("\tSender: %v\n\tReceiver: %v\n\tAmount: %d\n\tSignature: %x\n", tx1.Sender.Name, tx1.Receiver.Name, tx1.Amount, tx1.Signature)
			communalLedger.Print()
		}
	}
	fmt.Println("Simulating transaction tampering and verification")
	fmt.Println("\nTampering with a transaction and attempting to re-verify")
	tamperedTransaction := &communalLedger.Transactions[0]
	originalAmount := tamperedTransaction.Amount
	tamperedTransaction.Amount = 999

	fmt.Printf("Original amount: %d | Tampered amount: %d\n", originalAmount, tamperedTransaction.Amount)
	verificationResult := tamperedTransaction.Verify()
	fmt.Printf("Verification result after tampering: %v\n", verificationResult)
}
