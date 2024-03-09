package main

import (
	"fmt"

	"github.com/lf-hernandez/bitcoinplayground/ledger"
)

func main() {
	fmt.Println("Initializing new communal ledger")
	communalLedger := ledger.NewLedger()
	communalLedger.Print()

	fmt.Println("Distributing 20 coins to each participant")
	communalLedger.RunInitialDistribution()
	communalLedger.Print()

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
		fmt.Println("Felipe sent Satoshi 5 coin")
		err = communalLedger.SimulateTransaction(felipe, satoshi, 5)
		if err != nil {
			fmt.Println("Transaction error:  ", err)
		} else {
			communalLedger.Print()
		}
	}

	jane, err := communalLedger.GetParticipantByName("Jane")
	if err != nil {
		fmt.Println(err)
	}

	if jane != nil && satoshi != nil {
		fmt.Println("Satoshi sent Jane 15 coin")
		err = communalLedger.SimulateTransaction(satoshi, jane, 15)
		if err != nil {
			fmt.Println("Transaction error:  ", err)
		} else {
			communalLedger.Print()
		}
	}

}
