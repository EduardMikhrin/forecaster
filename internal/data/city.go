package data

type CitiesQ interface {
	GetAll() ([]City, error)
}

type City struct {
	Id   int64  `db:"id" structs:"-"`
	Name string `db:"name" structs:"name"`
}
