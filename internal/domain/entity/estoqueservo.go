package entity

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Estoqueservo struct {
	Submontagem_id int
	Db             *sql.DB `json:"-"`
}

func NovoEstoqueservo(db *sql.DB) (*Estoqueservo, error) {
	return &Estoqueservo{
		Submontagem_id: 0,
		Db:             db}, nil
}

func (es *Estoqueservo) CriarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	//Prgando os dados do Bady
	var estoqueservo Estoqueservo
	json.NewDecoder(r.Body).Decode(&estoqueservo)
	// inserindo no banco de dados
	query := "INSERT INTO safisa.estoqueservo (Submontagem_id) VALUES (?)"
	_, err := es.Db.Exec(query, estoqueservo.Submontagem_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
		return
	}

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Foi inserido no banco de dados corretamente")
}

func (es *Estoqueservo) ListarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	rows, err := es.Db.Query("SELECT * FROM safisa.estoqueservo")
	if err != nil {
		http.Error(w, "Falha ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var estoques []Estoqueservo
	for rows.Next() {
		var estoque Estoqueservo
		if err := rows.Scan(&estoque.Submontagem_id); err != nil {
			http.Error(w, "Falha ao ler os dados do banco de dados", http.StatusInternalServerError)
			return
		}
		estoques = append(estoques, estoque)
	}

	// Responda com a lista de estoques em formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estoques)

}

func (es *Estoqueservo) AtualizarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	// Extrair o ID da url
	id := chi.URLParam(r, "id")
	//convertendo Id do tipo string para int
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// pegando dados do Bary
	var estoque Estoqueservo
	err = json.NewDecoder(r.Body).Decode(&estoque)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding Peca: %v", err)
		return
	}

	// atualizando no banco de dados
	query := "UPDATE safisa.estoqueservo SET Submontagem_id =? WHERE Submontagem_id =?"
	_, err = es.Db.Exec(query, estoque.Submontagem_id, id_int)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
		return
	}

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Foi atualizado no banco de dados corretamente")
}

func (es *Estoqueservo) DeletarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	// extrauindo o Id da URL para excuir
	id := chi.URLParam(r, "id")
	//convertendo Id do tipo string para int
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//verificando se o id existe no banco de dados
	query := "SELECT * FROM safisa.estoqueservo WHERE Submontagem_id =?"
	row := es.Db.QueryRow(query, id_int)
	var estoque Estoqueservo
	err = row.Scan(&estoque.Submontagem_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
		return
	}

	// deletando no banco de dados
	query = "DELETE FROM safisa.estoqueservo WHERE Submontagem_id =?"
	_, err = es.Db.Exec(query, estoque.Submontagem_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("falha ao deletar dado do banco de dados")
		return
	}

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Foi deletado no banco de dados corretamente")
}

func (es *Estoqueservo) ListarEstoqueMpId(w http.ResponseWriter, r *http.Request) {

	// Extraindo Id da Url
	id := chi.URLParam(r, "id")
	//convertendo Id do tipo string para int
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//verificando se o id existe no banco de dados e ja mostrando na tela
	query := "SELECT * FROM safisa.estoqueservo WHERE Submontagem_id =?"
	row := es.Db.QueryRow(query, id_int)
	var estoque Estoqueservo
	err = row.Scan(&estoque.Submontagem_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
		return
	}

	// Responda com a lista de estoques em formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estoque)

}
