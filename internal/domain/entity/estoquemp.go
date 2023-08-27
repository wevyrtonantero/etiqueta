package entity

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

//se remete ao material para ser usinado, "materia prima"

type Estoquemp struct {
	Materiaprima_id int
	Db              *sql.DB `json:"-"`
}

func NovoEstoquemp(db *sql.DB) (*Estoquemp, error) {
	return &Estoquemp{
		Materiaprima_id: 0,
		Db:              db}, nil
}

func (emp *Estoquemp) CriarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	//Prgando os dados do Bady
	var estoquemp Estoquemp
	json.NewDecoder(r.Body).Decode(&estoquemp)
	// inserindo no banco de dados
	query := "INSERT INTO safisa.estoquemp (materiaprima_id) VALUES (?)"
	_, err := emp.Db.Exec(query, estoquemp.Materiaprima_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
		return
	}

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Foi inserido no banco de dados corretamente")
}

func (emp *Estoquemp) ListarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	rows, err := emp.Db.Query("SELECT * FROM safisa.estoquemp")
	if err != nil {
		http.Error(w, "Falha ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var estoques []Estoquemp
	for rows.Next() {
		var estoque Estoquemp
		if err := rows.Scan(&estoque.Materiaprima_id); err != nil {
			http.Error(w, "Falha ao ler os dados do banco de dados", http.StatusInternalServerError)
			return
		}
		estoques = append(estoques, estoque)
	}

	// Responda com a lista de estoques em formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estoques)

}

func (emp *Estoquemp) AtualizarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	// Extrair o ID da url
	id := chi.URLParam(r, "id")
	//convertendo Id do tipo string para int
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// pegando dados do Bary
	var estoque Estoquemp
	err = json.NewDecoder(r.Body).Decode(&estoque)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding Peca: %v", err)
		return
	}

	// atualizando no banco de dados
	query := "UPDATE safisa.estoquemp SET materiaprima_id =? WHERE materiaprima_id =?"
	_, err = emp.Db.Exec(query, estoque.Materiaprima_id, id_int)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
		return
	}

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Foi atualizado no banco de dados corretamente")
}

func (emp *Estoquemp) DeletarEstoqueMp(w http.ResponseWriter, r *http.Request) {
	// extrauindo o Id da URL para excuir 
    id := chi.URLParam(r, "id")
    //convertendo Id do tipo string para int
    id_int, err := strconv.Atoi(id)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	//verificando se o id existe no banco de dados	
	query := "SELECT * FROM safisa.estoquemp WHERE materiaprima_id =?"
    row := emp.Db.QueryRow(query, id_int)
    var estoque Estoquemp
    err = row.Scan(&estoque.Materiaprima_id)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
        return
    }

    // deletando no banco de dados
    query = "DELETE FROM safisa.estoquemp WHERE materiaprima_id =?"
    _, err = emp.Db.Exec(query, estoque.Materiaprima_id)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        json.NewEncoder(w).Encode("falha ao deletar dado do banco de dados")
        return
    }

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode("Foi deletado no banco de dados corretamente")
}

func (emp *Estoquemp) ListarEstoqueMpId(w http.ResponseWriter, r *http.Request) {

// Extraindo Id da Url
id := chi.URLParam(r, "id")
//convertendo Id do tipo string para int
id_int, err := strconv.Atoi(id)
if err!= nil {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}
//verificando se o id existe no banco de dados e ja mostrando na tela
query := "SELECT * FROM safisa.estoquemp WHERE materiaprima_id =?"
row := emp.Db.QueryRow(query, id_int)
var estoque Estoquemp
err = row.Scan(&estoque.Materiaprima_id)
if err!= nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
	return
}

// Responda com a lista de estoques em formato JSON
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(estoque)

}
