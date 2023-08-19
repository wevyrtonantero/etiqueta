package entity

import (
	"net/http"
)

type Pedido struct {
	Id             int
	Numerodepedido int
	Cliente        string
	Cidade         string
	Submontagem_id *SubMontagem
	Peca_id        *Peca
}

func CriarPedido(w http.ResponseWriter, r *http.Request) {

}

func ListarPedido(w http.ResponseWriter, r *http.Request) {

}

func AtualizarPedido(w http.ResponseWriter, r *http.Request) {

}
func DeletarPedido(w http.ResponseWriter, r *http.Request) {

}

func ListarPedidoid(w http.ResponseWriter, r *http.Request) {

}

func ListarPedidoId(w http.ResponseWriter, r *http.Request) {

}
