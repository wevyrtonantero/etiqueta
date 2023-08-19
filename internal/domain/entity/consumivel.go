package entity

import (
	"database/sql"
)

type Consumivel struct {
	Id           int
	Cod          string
	Descricao    string
	Categoria_id *Categoria
	imposto_id   *Imposto
	Db           *sql.DB `json:"-"`
}

func NovoConsumivel(db *sql.DB) (*Consumivel, error) {
	return &Consumivel{
		Id:           0,
		Cod:          "",
		Descricao:    "",
		Categoria_id: &Categoria{},
		imposto_id:   &Imposto{},
		Db:           db,
	}, nil
}
