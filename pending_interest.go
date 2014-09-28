package ndn

import "time"

type pendingInterest struct {
	ID       uint64
	Interest *Interest
	Data     chan *Data
	Timeout  chan time.Time
}

func (this *pendingInterest) deliver(d *Data) {
	this.Data <- d
	close(this.Data)
}
