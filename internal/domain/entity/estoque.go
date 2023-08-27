package entity

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Estoque struct {
	Peca_id int
	Db      *sql.DB `json:"-"`
}

func NovoEstoque(db *sql.DB) (*Estoque, error) {
	return &Estoque{
		Peca_id: 0,
		Db:      db,
	}, nil

}

func (e *Estoque) CriarEstoque(w http.ResponseWriter, r *http.Request) {
	// pegando dados do Bary
	var estoque Estoque
	err := json.NewDecoder(r.Body).Decode(&estoque)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding Peca: %v", err)
		return
	}

	// inserindo no banco de dados
	query := "INSERT INTO safisa.estoque (peca_id) VALUES (?)"
	_, err = e.Db.Exec(query, estoque.Peca_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
		return
	}

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Foi inserido no banco de dados corretamente")
}

func (e *Estoque) ListarEstoque(w http.ResponseWriter, r *http.Request) {
	//consultar estoque no banco de dados

	rows, err := e.Db.Query("SELECT * FROM safisa.estoque")
	if err != nil {
		http.Error(w, "Falha ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var estoques []Estoque
	for rows.Next() {
		var estoque Estoque
		if err := rows.Scan(&estoque.Peca_id); err != nil {
			http.Error(w, "Falha ao ler os dados do banco de dados", http.StatusInternalServerError)
			return
		}
		estoques = append(estoques, estoque)
	}

	// Responda com a lista de estoques em formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estoques)
}

func (e *Estoque) AtualizarEstoque(w http.ResponseWriter, r *http.Request) {
	// Extrair o ID da url
	id := chi.URLParam(r, "id")
	//convertendo Id do tipo string para int
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// pegando dados do Bary
	var estoque Estoque
    err = json.NewDecoder(r.Body).Decode(&estoque)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        log.Printf("Error decoding Peca: %v", err)
        return
    }

    // atualizando no banco de dados
    query := "UPDATE safisa.estoque SET peca_id =? WHERE peca_id =?"
    _, err = e.Db.Exec(query, estoque.Peca_id, id_int)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
        return
    }

    // retornando status 201
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode("Foi atualizado no banco de dados corretamente")
}

func (e *Estoque) DeletarEstoque(w http.ResponseWriter, r *http.Request) {
// extrauindo o Id da URL para excuir 
    id := chi.URLParam(r, "id")
    //convertendo Id do tipo string para int
    id_int, err := strconv.Atoi(id)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	//verificando se o id existe no banco de dados	
	query := "SELECT * FROM safisa.estoque WHERE peca_id =?"
    row := e.Db.QueryRow(query, id_int)
    var estoque Estoque
    err = row.Scan(&estoque.Peca_id)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
        return
    }

    // deletando no banco de dados
    query = "DELETE FROM safisa.estoque WHERE peca_id =?"
    _, err = e.Db.Exec(query, estoque.Peca_id)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        json.NewEncoder(w).Encode("falha ao deletar dado do banco de dados")
        return
    }

	// retornando status 201
	w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode("Foi deletado no banco de dados corretamente")

}

func (e *Estoque) ListarEstoqueId(w http.ResponseWriter, r *http.Request) {
// Extraindo Id da Url
    id := chi.URLParam(r, "id")
    //convertendo Id do tipo string para int
    id_int, err := strconv.Atoi(id)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	//verificando se o id existe no banco de dados
	query := "SELECT * FROM safisa.estoque WHERE peca_id =?"
    row := e.Db.QueryRow(query, id_int)
    var estoque Estoque
    err = row.Scan(&estoque.Peca_id)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        json.NewEncoder(w).Encode("Verifique se o id da peça existe ")
        return
    }

    // Responda com a lista de estoques em formato JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(estoque)


}
