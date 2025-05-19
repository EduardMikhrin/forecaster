package pg

import (
	"database/sql"
	"fmt"
	"github.com/EduardMikhrin/forecaster/internal/data"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	citiesTableName = "cities"
)

type cityQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func (c cityQ) GetAll() ([]data.City, error) {
	var res []data.City

	err := c.db.Select(&res, c.sql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}
	return res, err
}

func NewCityQ(db *pgdb.DB) data.CitiesQ {
	return &cityQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", citiesTableName)),
		upd: sq.Update(citiesTableName),
	}
}
