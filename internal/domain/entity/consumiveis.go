package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Consumivel struct {
	Id         int
	Cod        string
	Descricao  string
	Fornecedor []int
	Categoria  []int
	Db         *sql.DB `json:"-"`
}

func NovoConsumivel(db *sql.DB) (*Consumivel, error) {
	return &Consumivel{
		Id:         0,
		Cod:        "",
		Descricao:  "",
		Fornecedor: []int{},
		Categoria:  []int{},
		Db:         db,
	}, nil
}

func (cl *Consumivel) CriarConsumivel(w http.ResponseWriter, r *http.Request) {
	var consumivel Consumivel
	err := json.NewDecoder(r.Body).Decode(&consumivel)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}

	if consumivel.Cod == "" || len(consumivel.Fornecedor) == 0 || len(consumivel.Categoria) == 0 {
		http.Error(w, "Campos obrigatórios não preenchidos", http.StatusBadRequest)
		return
	}

	
	// Criar o Consumivel
	query := "INSERT INTO safisa.Consumivel (cod, descricao) VALUES (?, ?)"
	_, err = cl.Db.Exec(query, consumivel.Cod, consumivel.Descricao)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Consumível inserido com sucesso!")
}

func (cl *Consumivel) ListarConsumivel(w http.ResponseWriter, r *http.Request) {

}

func (cl *Consumivel) AtualizarConsumivel(w http.ResponseWriter, r *http.Request) {

}
func (cl *Consumivel) DeletarConsumivel(w http.ResponseWriter, r *http.Request) {

}

func (cl *Consumivel) ListarConsumivelId(w http.ResponseWriter, r *http.Request) {

}
