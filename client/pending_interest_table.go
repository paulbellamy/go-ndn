package client

type pendingInterestTable struct {
	nextID uint64
	items  map[uint64]*pendingInterest
}

func newPendingInterestTable() *pendingInterestTable {
	return &pendingInterestTable{
		nextID: 0,
		items:  make(map[uint64]*pendingInterest),
	}
}

func (p *pendingInterestTable) AddInterest(i *Interest) *pendingInterest {
	return nil
}
