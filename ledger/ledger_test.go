package ledger

import (
	"testing"
)

func Test_validateTransaction(t *testing.T) {
	sendToSelfPerson := &Person{Name: "Jane Smith", Balance: 40}
	testCases := []struct {
		name        string
		sender      *Person
		receiver    *Person
		amount      int
		shouldThrow bool
	}{
		{
			name:        "insufficient funds",
			sender:      &Person{Name: "John Smith", Balance: 10},
			receiver:    &Person{Name: "Jane Smith", Balance: 5},
			amount:      20,
			shouldThrow: true,
		},
		{
			name:        "negative amount",
			sender:      &Person{Name: "John Smith", Balance: 20},
			receiver:    &Person{Name: "Jane Smith", Balance: 20},
			amount:      -20,
			shouldThrow: true,
		},
		{
			name:        "sending to self",
			sender:      sendToSelfPerson,
			receiver:    sendToSelfPerson,
			amount:      10,
			shouldThrow: true,
		},
		{
			name:        "valid transaction",
			sender:      &Person{Name: "John Smith", Balance: 20},
			receiver:    &Person{Name: "Jane Smith", Balance: 20},
			amount:      10,
			shouldThrow: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validateTransaction(tc.sender, tc.receiver, tc.amount)
			if (result != nil) != tc.shouldThrow {
				t.Errorf("%v - incorrect results, expected %v, got %v", tc.name, tc.shouldThrow, result)
			}
		})
	}
}

func Test_addTransaction(t *testing.T) {
	ledger := NewLedger()

	john := ledger.Participants[0]
	jane := ledger.Participants[1]

	testCases := []struct {
		name         string
		sender       *Person
		receiver     *Person
		amount       int
		expectedSize int
	}{
		{
			name:         "Valid Transaction",
			sender:       john,
			receiver:     jane,
			amount:       10,
			expectedSize: len(ledger.Transactions) + 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ledger.addTransaction(tc.sender, tc.receiver, tc.amount)

			if len(ledger.Transactions) != tc.expectedSize {
				t.Errorf("Expected %d transactions, found %d", tc.expectedSize, len(ledger.Transactions))
				return
			}

			lastTransaction := ledger.Transactions[len(ledger.Transactions)-1]
			if lastTransaction.Sender != tc.sender || lastTransaction.Receiver != tc.receiver || lastTransaction.Amount != tc.amount {
				t.Errorf("Transaction was not added correctly")
			}
		})
	}
}
