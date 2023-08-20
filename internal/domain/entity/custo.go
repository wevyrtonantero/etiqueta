package entity

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Custo struct {
	Id        int
	Cod       string
	Descricao string
	Hmaquina  float64
	Zincagem  float64
	Retifica  float64
	Tempera   float64
	Dobra     float64
	Db        *sql.DB `json:"-"`
}

func NovoCusto(db *sql.DB) (*Custo, error) {
	return &Custo{
		Id:        0,
		Cod:       "",
		Descricao: "",
		Hmaquina:  0,
		Zincagem:  0,
		Retifica:  0,
		Tempera:   0,
		Dobra:     0,
		Db:        db,
	}, nil
}

func (cto *Custo) CriarCusto(w http.ResponseWriter, r *http.Request) {
	var custo Custo
	err := json.NewDecoder(r.Body).Decode(&custo)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}

	if custo.Cod == "" || custo.Descricao == "" {
		http.Error(w, "Campos 'Cod' e 'Descricao' são obrigatórios", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO safisa.custo (id, cod, descricao, hmaquina, zincagem, retifica, tempera, dobra) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = cto.Db.Exec(query, custo.Id, custo.Cod, custo.Descricao, custo.Hmaquina, custo.Zincagem, custo.Retifica, custo.Tempera, custo.Dobra)
	if err != nil {
		// Exiba o erro corretamente
		println("Erro ao inserir dados no banco de dados:", err.Error())
		http.Error(w, "Erro ao inserir dados no banco de dados", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cto *Custo) ListarCusto(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, cod, descricao, hmaquina, zincagem, retifica, tempera, dobra FROM safisa.custo"
	rows, err := cto.Db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao buscar dados do banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var custos []Custo
	for rows.Next() {
		var custo Custo
		err := rows.Scan(&custo.Id, &custo.Cod, &custo.Descricao, &custo.Hmaquina, &custo.Zincagem, &custo.Retifica, &custo.Tempera, &custo.Dobra)
		if err != nil {
			http.Error(w, "Erro ao ler dados do banco de dados", http.StatusInternalServerError)
			return
		}
		custos = append(custos, custo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(custos)
}

func (cto *Custo) AtualizarCusto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var novoCusto Custo
	err = json.NewDecoder(r.Body).Decode(&novoCusto)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}

	if novoCusto.Cod == "" || novoCusto.Descricao == "" {
		http.Error(w, "Campos 'Cod' e 'Descricao' são obrigatórios", http.StatusBadRequest)
		return
	}

	query := "UPDATE safisa.custo SET cod = ?, descricao = ?, hmaquina = ?, zincagem = ?, retifica = ?, tempera = ?, dobra = ? WHERE id = ?"
	result, err := cto.Db.Exec(query, novoCusto.Cod, novoCusto.Descricao, novoCusto.Hmaquina, novoCusto.Zincagem, novoCusto.Retifica, novoCusto.Tempera, novoCusto.Dobra, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar dados no banco de dados", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Erro ao obter número de linhas afetadas", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "ID não encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cto *Custo) DeletarCusto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM safisa.custo WHERE id = ?"
	_, err = cto.Db.Exec(query, id)
	if err != nil {
		http.Error(w, "Erro ao excluir dados do banco de dados", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (cto *Custo) ListarCustoID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	query := "SELECT id, cod, descricao, hmaquina, zincagem, retifica, tempera, dobra FROM safisa.custo WHERE id = ?"
	row := cto.Db.QueryRow(query, id)

	var custo Custo
	err = row.Scan(&custo.Id, &custo.Cod, &custo.Descricao, &custo.Hmaquina, &custo.Zincagem, &custo.Retifica, &custo.Tempera, &custo.Dobra)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Custo não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar dados do banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(custo)
}
