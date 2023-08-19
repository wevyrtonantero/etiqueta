package entity

import (
	"errors"
	"net/http"
)

type Fornecedor struct {
	Id       int
	Nome     string
	Telefone string
	Rua      string
	Numero   string
	Cidade   string
	Uf       string
	Contato  string
	
}

func (f *Fornecedor) IsValid() error {
	if f.Nome == "" {
		return errors.New("NOME_OBRIGATORIO")
	}
	return nil
}

func CriarFornecedor(w http.ResponseWriter, r *http.Request) {

}

func ListarFornecedor(w http.ResponseWriter, r *http.Request) {

}

func AtualizarFornecedor(w http.ResponseWriter, r *http.Request) {

}
func DeletarFornecedor(w http.ResponseWriter, r *http.Request) {

}
func ListarFornecedorId(w http.ResponseWriter, r *http.Request) {

}
