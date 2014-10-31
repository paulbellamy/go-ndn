package ndn

import (
	"time"

	"github.com/paulbellamy/go-ndn/packets"
)

type pendingInterest struct {
	ID       uint64
	Interest *packets.Interest
	Data     chan *packets.Data
	Timeout  chan time.Time
}

func (this *pendingInterest) deliver(d *packets.Data) {
	this.Data <- d
	close(this.Data)
}
