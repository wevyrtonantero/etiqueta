package entity

import (
	"errors"
	"net/http"
	"time"
)

type Pedido struct {
	NumeroPedido  int
	Cliente       string
	Cidade        string
	Peca          []Peca
	NumeroDeSerie []string
	Data          time.Time
	Datasaida     string
}

func (pd *Pedido) IsValid() error {
	if pd.NumeroPedido < 0 {
		return errors.New("CODIGO_OBRIGATORIO")
	}
	if pd.Cliente == "" {
		return errors.New("DESCRIÇAO_OBRIGATORIO")
	}

	return nil

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
