package ndn

type pendingInterest struct {
	ID       uint64
	Interest *Interest
	Data     chan *Data
}