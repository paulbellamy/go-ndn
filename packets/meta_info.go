package packets

import "time"

type MetaInfo struct {
	freshnessPeriod    time.Duration
	hasFreshnessPeriod bool
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
