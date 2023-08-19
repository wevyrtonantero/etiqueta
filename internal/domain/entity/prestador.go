package entity


type PrestadorDeServico struct {
	Id       int
	Nome     string
	Telefone int
	Rua      string
	Numero   string
	Cidade   string
	Uf       string
	Contato  string
	Peca     []int 
}