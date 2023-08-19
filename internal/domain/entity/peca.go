package entity

import (
	"errors"
	"net/http"
)

type Peca struct {
	Id                int      //Ind
	Codigo            string   // codigo interno safisa ex. A-011/12 - 071 - B-011/20
	Descricao         string   //Nome da peça ex. Niple de Ar
	Massa             float64  // Massa da Peça em Kg
	Consumomes        int      //quantidade consumida em 1 mes
	UrlFoto           string   // link Para acessar uma foto da peça
	UrlDesenho        string   // link para acessar o desenho técnico da peça
	UrlSetupDeMaquina string   // Link para acessar como se produz essa peça (setup de maquina)
	DescricaoTecnica  string   // Um breve resumo de onde a peça vai como é usado e suas caracteristicas tecnicas
	MateriaPrima      []int    // Que tipo de Materia Prima é usado nesta peça ex: Aço - Chumbaloy - 7/8'
	Fornecedor        []int    // onde compra ou qm faz
	Categoria         []int    //Qual a caregoria da peça - ex: Cotovelo - empurradores - conector de óleo

}

func (p *Peca) IsValid() error {
	if p.Codigo == "" {
		return errors.New("CODIGO_OBRIGATORIO")
	}
	if p.Descricao == "" {
		return errors.New("DESCRIÇAO_OBRIGATORIO")
	}
	return nil

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
