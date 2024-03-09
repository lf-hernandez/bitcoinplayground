package ledger

import (
	"errors"
	"fmt"
)

type Person struct {
	Name    string
	Balance int
}

func (p Person) String() string {
	return fmt.Sprintf("Name: %s, Balance: %d", p.Name, p.Balance)
}

type Transaction struct {
	Sender   Person
	Receiver Person
	Amount   int
}

func (t Transaction) String() string {
	return fmt.Sprintf("Sender: %s, Receiver: %s, Amount: %d", t.Sender.Name, t.Receiver.Name, t.Amount)
}

type Pot struct {
	Total int
}

func (p Pot) String() string {
	return fmt.Sprintf("%d", p.Total)
}

type Ledger struct {
	Pot          Pot
	Transactions []Transaction
	Participants []*Person
}

func NewLedger() *Ledger {
	return &Ledger{
		Pot:          Pot{Total: 100},
		Transactions: []Transaction{},
		Participants: []*Person{
			{Name: "Felipe", Balance: 0},
			{Name: "John", Balance: 0},
			{Name: "Satoshi", Balance: 0},
			{Name: "Jane", Balance: 0},
		},
	}
}

func (l *Ledger) RunInitialDistribution() {
	distAmount := 20
	for _, p := range l.Participants {
		l.Pot.Total -= distAmount
		p.Balance += distAmount
	}
}

func (l *Ledger) addTransaction(sender *Person, receiver *Person, amount int) {
	newTransaction := Transaction{Sender: *sender, Receiver: *receiver, Amount: amount}
	l.Transactions = append(l.Transactions, newTransaction)
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

func (l *Ledger) SimulateTransaction(sender *Person, receiver *Person, amount int) error {
	err := validateTransaction(sender, receiver, amount)
	if err != nil {
		return err
	}

	l.addTransaction(sender, receiver, amount)
	l.updateBalancePostTransaction(sender, receiver, amount)
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
		fmt.Print("\t\t[")
		for idx, t := range l.Transactions {
			fmt.Printf("\n\t\t%d - %v,", idx, t)
		}
		fmt.Println("\n\t\t]")
	}
}
