package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
	Db       *sql.DB `json:"-"`
}

func NovoFornecedor(db *sql.DB) (*Fornecedor, error) {
	return &Fornecedor{
		Id:       0,
		Nome:     "",
		Telefone: "",
		Rua:      "",
		Numero:   "",
		Cidade:   "",
		Uf:       "",
		Contato:  "",
		Db:       db,
	}, nil
}

func (f *Fornecedor) CriarFornecedor(w http.ResponseWriter, r *http.Request) {
	var fornecedor Fornecedor
	err := json.NewDecoder(r.Body).Decode(&fornecedor)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		log.Println("erro ao decodificar o json", err)
		return
	}

	if fornecedor.Id <= 0 || fornecedor.Nome == "" || fornecedor.Telefone == "" {
		http.Error(w, "Os campos id, nome e telefone são obrigatórios", http.StatusBadRequest)
		log.Println("Campo Id, Nome e Telefone são obrigatórios")
		return
	}

	query := "INSERT INTO fornecedor (id, nome, telefone, rua, numero, cidade, uf, contato) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = f.Db.Exec(query, fornecedor.Id, fornecedor.Nome, fornecedor.Telefone, fornecedor.Rua, fornecedor.Numero, fornecedor.Cidade, fornecedor.Uf, fornecedor.Contato)
	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco de dados", http.StatusInternalServerError)
		log.Println("Erro ao inserir dados no banco de dados:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Fornecedor criado com sucesso"))
}

func (f *Fornecedor) ListarFornecedor(w http.ResponseWriter, r *http.Request) {

	// Cria a consulta SQL para buscar todos os fornecedores
	query := `
		SELECT id, nome, telefone, rua, numero, cidade, uf, contato
		FROM safisa.fornecedor
	`

	// Executa a consulta SQL
	rows, err := f.Db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao buscar dados no banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Cria uma slice para armazenar os fornecedores
	var fornecedores []Fornecedor

	// Lê cada linha do resultado da consulta e cria um Fornecedor
	for rows.Next() {
		var fornecedor Fornecedor
		err := rows.Scan(&fornecedor.Id, &fornecedor.Nome, &fornecedor.Telefone, &fornecedor.Rua, &fornecedor.Numero, &fornecedor.Cidade, &fornecedor.Uf, &fornecedor.Contato)
		if err != nil {
			http.Error(w, "Erro ao ler dados do banco de dados", http.StatusInternalServerError)
			return
		}
		fornecedores = append(fornecedores, fornecedor)
	}

	// Retorna a lista de fornecedores como resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fornecedores)

}

func (f *Fornecedor) AtualizarFornecedor(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	var fornecedor Fornecedor
	err = json.NewDecoder(r.Body).Decode(&fornecedor)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Valida os campos obrigatórios
	if fornecedor.Nome == "" || fornecedor.Uf == "" {
		http.Error(w, "Campos 'Nome' e 'Uf' são obrigatórios", http.StatusBadRequest)
		return
	}

	// Consulta o banco de dados para verificar se o Fornecedor com o ID existe
	var count int
	err = f.Db.QueryRow("SELECT COUNT(*) FROM safisa.fornecedor WHERE id = ?", id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar o ID no banco de dados", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Fornecedor com ID especificado não encontrado", http.StatusNotFound)
		return
	}

	// Cria a consulta SQL para atualizar o Fornecedor por ID
	query := `
			UPDATE safisa.fornecedor
			SET nome = ?, telefone = ?, rua = ?, numero = ?, cidade = ?, uf = ?, contato = ?
			WHERE id = ?
		`

	// Executa a consulta SQL
	_, err = f.Db.Exec(query, fornecedor.Nome, fornecedor.Telefone, fornecedor.Rua, fornecedor.Numero, fornecedor.Cidade, fornecedor.Uf, fornecedor.Contato, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Fornecedor atualizado com sucesso!")

}
func (f *Fornecedor) DeletarFornecedor(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	// Consulta o banco de dados para verificar se o Fornecedor com o ID existe
	var count int
	err = f.Db.QueryRow("SELECT COUNT(*) FROM safisa.fornecedor WHERE id = ?", id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar o ID no banco de dados", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Fornecedor com ID especificado não encontrado", http.StatusNotFound)
		return
	}

	// Cria a consulta SQL para excluir o Fornecedor por ID
	query := "DELETE FROM safisa.fornecedor WHERE id = ?"
	_, err = f.Db.Exec(query, id)
	if err != nil {
		http.Error(w, "Erro ao excluir dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Fornecedor excluído com sucesso!")

}
func (f *Fornecedor) ListarFornecedorId(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	// Consulta o banco de dados para buscar o Fornecedor por ID
	query := `
	SELECT id, nome, telefone, rua, numero, cidade, uf, contato
	FROM safisa.fornecedor
	WHERE id = ?
`
	row := f.Db.QueryRow(query, id)

	var fornecedor Fornecedor
	err = row.Scan(&fornecedor.Id, &fornecedor.Nome, &fornecedor.Telefone, &fornecedor.Rua, &fornecedor.Numero, &fornecedor.Cidade, &fornecedor.Uf, &fornecedor.Contato)
	if err != nil {
		http.Error(w, "Erro ao buscar o Fornecedor no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna o Fornecedor encontrado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fornecedor)
}
