package pg

import (
	"database/sql"
	"github.com/EduardMikhrin/forecaster/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type masterQ struct {
	db *pgdb.DB
}

func (m masterQ) New() data.MasterQ {
	return &masterQ{
		db: m.db,
	}
}

func (m masterQ) CitiesQ() data.CitiesQ {
	return NewCityQ(m.db)
}

func (m masterQ) Transaction(fn func(data interface{}) error, i interface{}) error {
	return m.db.TransactionWithOptions(&sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}, func() error {
		return fn(i)
	})
}

func (m masterQ) SubscriptionQ() data.SubscriptionQ {
	return NewSubscriptionQ(m.db)
}
