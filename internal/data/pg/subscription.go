package pg

import (
	"database/sql"
	"fmt"
	"github.com/EduardMikhrin/forecaster/internal/data"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	subscriptionsTableName = "subscriptions"
	subscriptionEmailField = "email"
)

type subscriptionQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func (s subscriptionQ) New() data.SubscriptionQ {
	return NewSubscriptionQ(s.db.Clone())
}

func (s subscriptionQ) Get() (*data.Subscription, error) {

	res := data.Subscription{}
	s.sql = sq.Select(
		"s.email",
		"s.city_id",
		"c.name AS city",
		"s.created_at",
	).PlaceholderFormat(sq.Dollar)

	err := s.db.Get(&res, s.sql.
		From("subscriptions AS s").
		LeftJoin("cities AS c ON s.city_id = c.id"),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, data.ErrNotFound
		}

		return nil, errors.Wrap(err, "error getting subscription")
	}

	return &res, nil
}

func (s subscriptionQ) GetAll() ([]data.Subscription, error) {
	res := []data.Subscription{}

	err := s.db.Select(&res, s.sql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, data.ErrNotFound
		}

		return nil, errors.Wrap(err, "error getting subscriptions")
	}

	return res, nil
}

func (s subscriptionQ) Insert(data *data.Subscription) error {
	clauses := structs.Map(data)

	if err := s.db.Exec(sq.Insert(subscriptionsTableName).SetMap(clauses)); err != nil {
		if pgdb.IsConstraintErr(err, subscriptionEmailField) {
			return data.ErrAlreadyExists
		}
		return errors.Wrap(err, "failed to insert subscription")
	}

	return nil
}

func (s subscriptionQ) Delete(email string) error {

	if err := s.db.Exec(sq.Delete(subscriptionsTableName).Where(sq.Eq{subscriptionEmailField: email})); err != nil {
		return errors.Wrap(err, "failed to delete subscription")
	}

	return nil
}

func (s subscriptionQ) FilterByEmail(email string) data.SubscriptionQ {
	s.sql = s.sql.Where(sq.Eq{subscriptionEmailField: email})
	s.upd = s.upd.Where(sq.Eq{subscriptionEmailField: email})

	return s
}

func NewSubscriptionQ(db *pgdb.DB) data.SubscriptionQ {
	return &subscriptionQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", subscriptionsTableName)),
		upd: sq.Update(subscriptionsTableName),
	}
}
