package scripts

/*package scripts

import (
	"github.com/adriel-meb/appointly-backend/internal/models"
	"time"
)

func (p *Provider) ComputeNextAvailable() {
	now := time.Now()
	var nextSlot *Availability

	for _, slot := range p.Availabilities {
		if slot.Date != nil && slot.Date.After(now) {
			if nextSlot == nil || slot.Date.Before(*nextSlot.Date) {
				nextSlot = &slot
			}
		}
	}

	if nextSlot != nil && nextSlot.Date != nil {
		diff := nextSlot.Date.Sub(now)
		if diff.Hours() < 24 {
			p.NextAvailable = "Aujourd'hui"
		} else if diff.Hours() < 48 {
			p.NextAvailable = "Demain"
		} else {
			p.NextAvailable = nextSlot.Date.Format("02 Jan 2006")
		}
	} else {
		p.NextAvailable = "Aucune disponibilité"
	}
}
*/
