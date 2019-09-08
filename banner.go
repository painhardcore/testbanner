package banner

import (
	"errors"
	"sort"

	"github.com/m-rec/06500a2490eaf9f133be55a8cb7c01c5ab0c9d45/ip"
)

// DisplayBanner is the banner representation
type DisplayBanner struct {
	Promotion
	img string
	src string
}

type banner struct {
	srs []DisplayBanner
}

// New sorts slice of DisplayBanner and provides APi for it.
func New(srs []DisplayBanner) (*banner, error) {
	// if there is more than 1 banner - presort it.
	switch len(srs) {
	case 0:
		return nil, errors.New("empty banners list")
	case 1:
	default:
		sort.Slice(srs, func(i, j int) bool {
			return srs[i].Expiration().Before(srs[j].Expiration())
		})
	}
	return &banner{srs: srs}, nil
}

// DisplayFor is used to get desired banner which fits specifications
// if you got nil - there are no available banner to use
func (b *banner) DisplayFor(addr string) *DisplayBanner {
	if ip.IsInternal(addr) {
		return b.firstNonExpired()
	}
	return b.firstActive()
}

// firstActive grabs the first active banner
func (b *banner) firstActive() *DisplayBanner {
	for i := range b.srs {
		if b.srs[i].Active() {
			return &b.srs[i]
		}
	}
	return nil
}

// firstActive grabs the first non expired one
func (b *banner) firstNonExpired() *DisplayBanner {
	for i := range b.srs {
		if !b.srs[i].Expired() {
			return &b.srs[i]
		}
	}
	return nil
}
