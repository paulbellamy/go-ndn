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
	p.nextID++
	pi := &pendingInterest{
		ID:       p.nextID,
		Interest: i,
	}
	p.items[pi.ID] = pi
	return pi
}

func (p *pendingInterestTable) RemovePendingInterest(id uint64) {
	delete(p.items, id)
}

func (p *pendingInterestTable) DispatchData(d *Data) {
}
