package data

import "time"

type SubscriptionQ interface {
	New() SubscriptionQ
	Get() (*Subscription, error)
	GetAll() ([]Subscription, error)
	Insert(data *Subscription) error
	Delete(email string) error
	FilterByEmail(email string) SubscriptionQ
}

type Subscription struct {
	Email     string    `db:"email" structs:"email"`
	City      string    `db:"city" structs:"-"`
	CityId    uint64    `db:"city_id" structs:"city_id"`
	CreatedAt time.Time `db:"created_at" structs:"-"`
}
