package ledger

import "testing"

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
