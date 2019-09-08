package banner

import "time"

type Promotion struct {
	Name  string
	Start time.Time
	End   time.Time
}

func (p *Promotion) Expired() bool {
	return time.Now().After(p.End)
}

func (p *Promotion) Active() bool {
	return time.Now().After(p.Start) && time.Now().Before(p.End)
}

func (p *Promotion) Expiration() time.Time {
	return p.End
}
