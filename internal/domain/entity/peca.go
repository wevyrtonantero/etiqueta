package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
	Custopeca         float64
	Categoria_id      int
	Imposto_id        int
	Custo_id          int
	Materiaprima_id   int
	Db                *sql.DB
}

func NovaPeca(db *sql.DB) (*Peca, error) {
	return &Peca{
		Id:                0,
		Codigo:            "",
		Descricao:         "",
		Massa:             0,
		Consumomes:        0,
		UrlFoto:           "",
		UrlDesenho:        "",
		UrlSetupDeMaquina: "",
		DescricaoTecnica:  "",
		Custopeca:         0,
		Categoria_id:      0,
		Imposto_id:        0,
		Custo_id:          0,
		Materiaprima_id:   0,
		Db:                db,
	}, nil
}

func (p *Peca) CriarPeca(w http.ResponseWriter, r *http.Request) {
	// Decodifica o JSON do corpo da solicitação para uma nova instância de Peca
	var novaPeca Peca
	err := json.NewDecoder(r.Body).Decode(&novaPeca)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Validação dos campos obrigatórios
	if novaPeca.Codigo == "" || novaPeca.Descricao == "" || novaPeca.Consumomes == 0 {
		http.Error(w, "Campos 'codigo', 'descricao' e 'consumomes' são obrigatórios", http.StatusBadRequest)
		return
	}

	// Cria a consulta SQL para inserir a nova peça no banco de dados
	query := `
			INSERT INTO safisa.peca (id, codigo, descricao, consumomes, massa, urlfoto, urldesenho, urlsetupdemaquina, descricaotecnica, custopeca, categoria_id, imposto_id, custo_id, materiaprima_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

	// Executa a consulta SQL
	_, err = p.Db.Exec(query, novaPeca.Id, novaPeca.Codigo, novaPeca.Descricao, novaPeca.Consumomes, novaPeca.Massa, novaPeca.UrlFoto, novaPeca.UrlDesenho, novaPeca.UrlSetupDeMaquina, novaPeca.DescricaoTecnica, novaPeca.Custopeca, novaPeca.Categoria_id, novaPeca.Imposto_id, novaPeca.Custo_id, novaPeca.Materiaprima_id)
	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Peça criada com sucesso!")
}

func (p *Peca) ListarPeca(w http.ResponseWriter, r *http.Request) {
	// Cria a consulta SQL para listar todas as peças
	query := `
		SELECT id, codigo, descricao, consumomes, massa, urlfoto, urldesenho, urlsetupdemaquina, descricaotecnica, custopeca, categoria_id, imposto_id, custo_id, materiaprima_id
		FROM safisa.peca
	`

	// Executa a consulta SQL
	rows, err := p.Db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Cria um slice para armazenar as peças
	var pecas []Peca

	// Itera sobre os resultados da consulta e preenche o slice
	for rows.Next() {
		var peca Peca
		err := rows.Scan(&peca.Id, &peca.Codigo, &peca.Descricao, &peca.Consumomes, &peca.Massa, &peca.UrlFoto, &peca.UrlDesenho, &peca.UrlSetupDeMaquina, &peca.DescricaoTecnica, &peca.Custopeca, &peca.Categoria_id, &peca.Imposto_id, &peca.Custo_id, &peca.Materiaprima_id)
		if err != nil {
			http.Error(w, "Erro ao ler os dados do banco de dados", http.StatusInternalServerError)
			return
		}
		pecas = append(pecas, peca)
	}

	// Retorna a lista de peças em formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pecas)

}

func (p *Peca) AtualizarPeca(w http.ResponseWriter, r *http.Request) {

	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Decodifica o JSON do corpo da solicitação para uma nova instância de Peca
	var pecaAtualizada Peca
	err = json.NewDecoder(r.Body).Decode(&pecaAtualizada)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validação dos campos obrigatórios
	if pecaAtualizada.Codigo == "" || pecaAtualizada.Descricao == "" || pecaAtualizada.Consumomes <= 0 {
		http.Error(w, "Campos obrigatórios não preenchidos", http.StatusBadRequest)
		return
	}

	// Verifica se o ID existe no banco de dados
	queryExistencia := "SELECT COUNT(*) FROM safisa.peca WHERE id = ?"
	var countExistencia int
	err = p.Db.QueryRow(queryExistencia, id).Scan(&countExistencia)
	if err != nil {
		http.Error(w, "Erro ao verificar a existência do ID", http.StatusInternalServerError)
		return
	}
	if countExistencia == 0 {
		http.Error(w, "ID não encontrado", http.StatusNotFound)
		return
	}

	// Cria a consulta SQL para atualizar a peça no banco de dados
	queryAtualizacao := `
	UPDATE safisa.peca
	SET codigo=?, descricao=?, consumomes=?, massa=?, urlfoto=?, urldesenho=?, urlsetupdemaquina=?, descricaotecnica=?, custopeca=?, categoria_id=?, imposto_id=?, custo_id=?, materiaprima_id=?
	WHERE id=?
`

	// Executa a consulta SQL de atualização
	_, err = p.Db.Exec(queryAtualizacao, pecaAtualizada.Codigo, pecaAtualizada.Descricao, pecaAtualizada.Consumomes, pecaAtualizada.Massa, pecaAtualizada.UrlFoto, pecaAtualizada.UrlDesenho, pecaAtualizada.UrlSetupDeMaquina, pecaAtualizada.DescricaoTecnica, pecaAtualizada.Custopeca, pecaAtualizada.Categoria_id, pecaAtualizada.Imposto_id, pecaAtualizada.Custo_id, pecaAtualizada.Materiaprima_id, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Peça atualizada com sucesso!")

}

func (p *Peca) DeletarPeca(w http.ResponseWriter, r *http.Request) {

// Extrai o ID do parâmetro da URL
idStr := chi.URLParam(r, "id")
id, err := strconv.Atoi(idStr)
if err != nil {
	http.Error(w, "ID inválido", http.StatusBadRequest)
	return
}

// Verifica se o ID existe no banco de dados
queryExistencia := "SELECT COUNT(*) FROM safisa.peca WHERE id = ?"
var countExistencia int
err = p.Db.QueryRow(queryExistencia, id).Scan(&countExistencia)
if err != nil {
	http.Error(w, "Erro ao verificar a existência do ID", http.StatusInternalServerError)
	return
}
if countExistencia == 0 {
	http.Error(w, "ID não encontrado", http.StatusNotFound)
	return
}

// Cria a consulta SQL para deletar a peça do banco de dados
queryDelecao := "DELETE FROM safisa.peca WHERE id = ?"

// Executa a consulta SQL de deleção
_, err = p.Db.Exec(queryDelecao, id)
if err != nil {
	http.Error(w, "Erro ao deletar a peça do banco de dados", http.StatusInternalServerError)
	return
}

// Retorna uma resposta de sucesso
w.WriteHeader(http.StatusOK)
fmt.Fprintln(w, "Peça deletada com sucesso!")



}

func (p *Peca) ListarPecaId(w http.ResponseWriter, r *http.Request) {
// Extrai o ID do parâmetro da URL
idStr := chi.URLParam(r, "id")
id, err := strconv.Atoi(idStr)
if err != nil {
	http.Error(w, "ID inválido", http.StatusBadRequest)
	return
}

// Cria a consulta SQL para buscar a peça por ID
query := "SELECT id, codigo, descricao, massa, consumomes, urlfoto, urldesenho, urlsetupdemaquina, descricaotecnica, custopeca, categoria_id, imposto_id, custo_id, materiaprima_id FROM safisa.peca WHERE id = ?"

// Executa a consulta SQL
row := p.Db.QueryRow(query, id)

// Cria uma instância de Peca para armazenar os resultados
var peca Peca
err = row.Scan(&peca.Id, &peca.Codigo, &peca.Descricao, &peca.Massa, &peca.Consumomes, &peca.UrlFoto, &peca.UrlDesenho, &peca.UrlSetupDeMaquina, &peca.DescricaoTecnica, &peca.Custopeca, &peca.Categoria_id, &peca.Imposto_id, &peca.Custo_id, &peca.Materiaprima_id)
if err != nil {
	http.Error(w, "Erro ao buscar a peça por ID", http.StatusInternalServerError)
	return
}

// Codifica a peça em JSON e escreve na resposta
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(peca)


}
