package entity

import (
	"errors"
	"net/http"
)

// aço - inox - aluminio -bronze
// chumbaloy - aluminio 3320 - inox 316l - bronze3350
type MateriaPrima struct {
	Id        int
	Descrição string
	Material  string
	Liga      string
	Kgmt      float64
	Valorkg   float64 // Massa kg por metro
	Peca      []string
	Imposto   []int
	Custo     []int
}

//lebrar de fazer a função calculando por barras e por metro

func (ma *MateriaPrima) IsValid() error {
	if ma.Descrição == "" {
		return errors.New("CODIGO_OBRIGATORIO")
	}
	if ma.Material == "" {
		return errors.New("CATEGORIA_OBRIGATORIO")
	}
	if ma.Liga == "" {
		return errors.New("CATEGORIA_OBRIGATORIO")
	}

	return nil
}
func CriarMateriaPrima(w http.ResponseWriter, r *http.Request) {

}

func ListarMateriaPrima(w http.ResponseWriter, r *http.Request) {

}

func AtualizarMateriaPrima(w http.ResponseWriter, r *http.Request) {

}
func DeletarMateriaPrima(w http.ResponseWriter, r *http.Request) {

}

func ListarMateriaPrimaId(w http.ResponseWriter, r *http.Request) {

}
