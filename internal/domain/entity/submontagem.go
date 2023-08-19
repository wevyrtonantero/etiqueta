package entity

import (
	
	"net/http"
)

type SubMontagem struct {
	Id int
	Cod       string
	Descricao string
	Urlfoto string
	Peca_id *Peca
}





func CriarSubMontagem(w http.ResponseWriter, r *http.Request) {

}

func ListarSubMontagem(w http.ResponseWriter, r *http.Request) {

}

func AtualizarSubMontagem(w http.ResponseWriter, r *http.Request) {

}
func DeletarSubMontagem(w http.ResponseWriter, r *http.Request) {

}
func ListarSubMontagemId(w http.ResponseWriter, r *http.Request) {

}
