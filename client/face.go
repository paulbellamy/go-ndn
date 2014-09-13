package client

func Face(transport Transport) *face {
	return &face{
		transport:            transport,
		pendingInterestTable: newPendingInterestTable(),
	}
}

type face struct {
	transport            Transport
	pendingInterestTable *pendingInterestTable
}

func (f *face) ExpressInterest(i *Interest) (*pendingInterest, error) {
	pendingInterest := f.pendingInterestTable.AddInterest(i)
	_, err := i.WriteTo(f.transport)
	return pendingInterest, err
}
