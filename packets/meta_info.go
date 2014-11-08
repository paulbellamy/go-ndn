package packets

import (
	"time"

	"github.com/paulbellamy/go-ndn/name"
)

const (
	UnknownContentType ContentType = -1
	DefaultContentType ContentType = 0
	LinkContentType    ContentType = 1
	KeyContentType     ContentType = 2
)

type ContentType int

type MetaInfo struct {
	contentType    ContentType
	hasContentType bool

	freshnessPeriod    time.Duration
	hasFreshnessPeriod bool

	finalBlockID name.Component
}

func (m *MetaInfo) GetContentType() ContentType {
	if !m.hasContentType {
		return UnknownContentType
	}
	return m.contentType
}

func (m *MetaInfo) SetContentType(x ContentType) {
	m.hasContentType = true
	m.contentType = x
}

func (m *MetaInfo) GetFreshnessPeriod() time.Duration {
	if !m.hasFreshnessPeriod {
		return -1
	}
	return m.freshnessPeriod
}

func (m *MetaInfo) SetFreshnessPeriod(x time.Duration) {
	m.hasFreshnessPeriod = true
	m.freshnessPeriod = x
}

func (m *MetaInfo) GetFinalBlockID() name.Component {
	return m.finalBlockID
}

func (m *MetaInfo) SetFinalBlockID(x name.Component) {
	m.finalBlockID = x
}
