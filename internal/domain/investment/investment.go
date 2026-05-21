package investment

import "time"

type Investment struct {
	ID        uint
	Name      string
	Balance   float64
	CreatedAt time.Time
}
