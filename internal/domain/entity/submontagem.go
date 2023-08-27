package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type SubMontagem struct {
	Id        int
	Cod       string
	Descricao string
	Urlfoto   string
	Peca_id   int
	Db        *sql.DB `json:"-"`
}

func NovaSubMontagem(db *sql.DB) (*SubMontagem, error) {
	return &SubMontagem{
		Id:        0,
		Cod:       "",
		Descricao: "",
		Urlfoto:   "",
		Peca_id:   0,
		Db:        db,
	}, nil
}

func (sb *SubMontagem) CriarSubMontagem(w http.ResponseWriter, r *http.Request) {
	//Pegando dados do Body
	var submontagem SubMontagem
	json.NewDecoder(r.Body).Decode(&submontagem)
	//validando os campos obrigatorios
	if submontagem.Cod == "" || submontagem.Descricao == "" || submontagem.Peca_id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Os campos obrigatórios não foram preenchidos!"))
		return
	}

	// Inserir a nova submontagem no banco de dados
	query := "INSERT INTO safisa.submontagem (id, cod, descricao, urlfoto, peca_id) VALUES (?, ?, ?, ?, ?)"
	_, err := sb.Db.Exec(query, submontagem.Id, submontagem.Cod, submontagem.Descricao, submontagem.Urlfoto, submontagem.Peca_id)
	if err != nil {
		fmt.Println("Erro ao inserir submontagem no banco de dados:", err) // Adicione esta linha
		http.Error(w, "Erro ao inserir submontagem no banco de dados, verifique se o id esta sendo digitado corretamente", http.StatusInternalServerError)
		return
	}

	// Responder com uma mensagem de sucesso
	response := map[string]string{"message": "Submontagem criada com sucesso"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (sb *SubMontagem) ListarSubMontagem(w http.ResponseWriter, r *http.Request) {
	// Consultar todas as submontagens no banco de dados
	query := "SELECT * FROM safisa.submontagem"
	rows, err := sb.Db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao consultar submontagens no banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Criar uma estrutura para armazenar as submontagens
	var submontagens []SubMontagem

	// Percorrer as linhas do resultado da consulta
	for rows.Next() {
		var submontagem SubMontagem
		err := rows.Scan(&submontagem.Id, &submontagem.Cod, &submontagem.Descricao, &submontagem.Urlfoto, &submontagem.Peca_id)
		if err != nil {
			http.Error(w, "Erro ao ler submontagem do resultado", http.StatusInternalServerError)
			return
		}
		submontagens = append(submontagens, submontagem)
	}

	// Responder com a lista de submontagens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submontagens)

}

func (sb *SubMontagem) AtualizarSubMontagem(w http.ResponseWriter, r *http.Request) {
	// Pegando o ID da submontagem a ser atualizada a partir dos parâmetros da URL usando o "chi"
	submontagemIDStr := chi.URLParam(r, "id")
	if submontagemIDStr == "" {
		http.Error(w, "ID da submontagem não fornecido", http.StatusBadRequest)
		return
	}

	// Convertendo o ID para um tipo inteiro
	submontagemID, err := strconv.Atoi(submontagemIDStr)
	if err != nil {
		http.Error(w, "ID da submontagem inválido", http.StatusBadRequest)
		return
	}

	// Consultar a submontagem no banco de dados pelo ID
	query := "SELECT * FROM safisa.submontagem WHERE id = ?"
	row := sb.Db.QueryRow(query, submontagemID)

	// Criar uma submontagem temporária para armazenar os dados existentes
	var submontagemExistente SubMontagem
	err = row.Scan(&submontagemExistente.Id, &submontagemExistente.Cod, &submontagemExistente.Descricao, &submontagemExistente.Urlfoto, &submontagemExistente.Peca_id)
	if err != nil {
		http.Error(w, "Submontagem não encontrada", http.StatusNotFound)
		return
	}

	// Decodificar o corpo JSON da requisição para atualizar a submontagem
	var submontagemAtualizada SubMontagem
	err = json.NewDecoder(r.Body).Decode(&submontagemAtualizada)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Atualizar a submontagem no banco de dados
	query = "UPDATE safisa.submontagem SET cod = ?, descricao = ?, urlfoto = ?, peca_id = ? WHERE id = ?"
	_, err = sb.Db.Exec(query, submontagemAtualizada.Cod, submontagemAtualizada.Descricao, submontagemAtualizada.Urlfoto, submontagemAtualizada.Peca_id, submontagemID)
	if err != nil {
		http.Error(w, "Erro ao atualizar submontagem no banco de dados", http.StatusInternalServerError)
		return
	}

	// Responder com uma mensagem de sucesso
	response := map[string]string{"message": "Submontagem atualizada com sucesso"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func (sb *SubMontagem) DeletarSubMontagem(w http.ResponseWriter, r *http.Request) {
	//pegando dados da URL para deletar
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID da submontagem não fornecido", http.StatusBadRequest)
		return
	}
	// Convertendo o ID para um tipo inteiro
	submontagemID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID da submontagem inválido", http.StatusBadRequest)
		return
	}

	// Consultar a submontagem no banco de dados pelo ID
	query := "SELECT * FROM safisa.submontagem WHERE id =?"
	row := sb.Db.QueryRow(query, submontagemID)
	// Criar uma submontagem temporária para armazenar os dados existentes
	var submontagemExistente SubMontagem
	err = row.Scan(&submontagemExistente.Id, &submontagemExistente.Cod, &submontagemExistente.Descricao, &submontagemExistente.Urlfoto, &submontagemExistente.Peca_id)
	if err != nil {
		http.Error(w, "Submontagem não encontrada", http.StatusNotFound)
		return
	}

	// Deletar a submontagem no banco de dados
	query = "DELETE FROM safisa.submontagem WHERE id =?"
	_, err = sb.Db.Exec(query, submontagemID)
	if err != nil {
		http.Error(w, "Erro ao deletar submontagem no banco de dados", http.StatusInternalServerError)
		return
	}
	// Responder com uma mensagem de sucesso
	response := map[string]string{"message": "Submontagem deletada com sucesso"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func (sb *SubMontagem) ListarSubMontagemId(w http.ResponseWriter, r *http.Request) {
	// Pegando o parâmetro "id" da URL usando o "chi"
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID da submontagem não fornecido", http.StatusBadRequest)
		return
	}

	// Convertendo o ID para um tipo inteiro
	submontagemID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID da submontagem inválido", http.StatusBadRequest)
		return
	}

	// Consultar a submontagem no banco de dados pelo ID
	query := "SELECT * FROM submontagem WHERE id = ?"
	row := sb.Db.QueryRow(query, submontagemID)

	// Criar uma submontagem temporária para armazenar os dados existentes
	var submontagemExistente SubMontagem
	err = row.Scan(&submontagemExistente.Id, &submontagemExistente.Cod, &submontagemExistente.Descricao, &submontagemExistente.Urlfoto, &submontagemExistente.Peca_id)
	if err != nil {
		http.Error(w, "Submontagem não encontrada", http.StatusNotFound)
		return
	}

	// Responder com a submontagem encontrada
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submontagemExistente)
}

