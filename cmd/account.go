package main

type Account struct {
	ActiveCard     bool `json:"activeCard"`
	AvailableLimit int  `json:"availableLimit"`
	transactions   []Transaction
}

func (acc *Account) countMatches(newTransaction Transaction) matches {
	matches := matches{}
	for _, t := range acc.transactions {
		minutesSinceLastTransaction := newTransaction.Time.Sub(t.Time).Minutes()
		if minutesSinceLastTransaction > IntervalMinutes {
			continue
		}
		if newTransaction.isSimilar(t) {
			matches.similarity++
		}
		matches.frequency++
	}
	return matches
}

type matches struct {
	frequency  int
	similarity int
}
