package entity

import (
	"net/http"
)

type Peca struct {
	Id                int     //Ind
	Codigo            string  // codigo interno safisa ex. A-011/12 - 071 - B-011/20
	Descricao         string  //Nome da peça ex. Niple de Ar
	Massa             float64 // Massa da Peça em Kg
	Consumomes        int     //quantidade consumida em 1 mes
	UrlFoto           string  // link Para acessar uma foto da peça
	UrlDesenho        string  // link para acessar o desenho técnico da peça
	UrlSetupDeMaquina string  // Link para acessar como se produz essa peça (setup de maquina)
	DescricaoTecnica  string  // Um breve resumo de onde a peça vai como é usado e suas caracteristicas tecnicas
	Categoria_id      *Categoria
	Imposto_id        *Imposto
	Custo_id          *Custo
	Materiaprima_id   *MateriaPrima
}

func CriarPeca(w http.ResponseWriter, r *http.Request) {

}

func ListarPeca(w http.ResponseWriter, r *http.Request) {

}

func AtualizarPeca(w http.ResponseWriter, r *http.Request) {

}
func DeletarPeca(w http.ResponseWriter, r *http.Request) {

}

func ListarPecaId(w http.ResponseWriter, r *http.Request) {

}
