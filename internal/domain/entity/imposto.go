package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Imposto struct {
	Id    int
	Uf    string
	Icms  float64
	Ipi   float64
	Difal float64
	Db    *sql.DB `json:"-"`
}

func NovoImposto(db *sql.DB) (*Imposto, error) {
	return &Imposto{
		Id:    0,
		Uf:    "",
		Icms:  0,
		Ipi:   0,
		Difal: 0,
		Db:    db,
	}, nil
}

func (i *Imposto) CriarImposto(w http.ResponseWriter, r *http.Request) {
	var imposto Imposto
	err := json.NewDecoder(r.Body).Decode(&imposto)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if imposto.Uf == "" {
		http.Error(w, "Campo 'uf' é obrigatório", http.StatusBadRequest)
		return
	}
	// Verificar se o ID já existe
	queryCheckID := "SELECT COUNT(*) FROM safisa.imposto WHERE id = ?"
	var count int
	err = i.Db.QueryRow(queryCheckID, imposto.Id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar a existência do ID", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "ID já existe", http.StatusBadRequest)
		return
	}

	// Continuar com a inserção no banco de dados

	query := "INSERT INTO safisa.imposto (id, uf, icms, ipi, difal) VALUES (?, ?, ?, ?, ?)"
	_, err = i.Db.Exec(query, imposto.Id, imposto.Uf, imposto.Icms, imposto.Ipi, imposto.Difal)
	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco de dados", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Imposto criado com sucesso!")
}

func (i *Imposto) ListarImposto(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, uf, icms, ipi, difal FROM safisa.imposto"
	rows, err := i.Db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao buscar dados do banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var impostos []Imposto
	for rows.Next() {
		var imposto Imposto
		err := rows.Scan(&imposto.Id, &imposto.Uf, &imposto.Icms, &imposto.Ipi, &imposto.Difal)
		if err != nil {
			http.Error(w, "Erro ao ler dados do banco de dados", http.StatusInternalServerError)
			return
		}
		impostos = append(impostos, imposto)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(impostos)
}

func (i *Imposto) AtualizarImposto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var novoImposto Imposto
	err = json.NewDecoder(r.Body).Decode(&novoImposto)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if novoImposto.Uf == "" {
		http.Error(w, "Campo 'uf' é obrigatório", http.StatusBadRequest)
		return
	}

	// Verificar se o ID já existe
	queryCheckID := "SELECT COUNT(*) FROM safisa.imposto WHERE id = ?"
	var count int
	err = i.Db.QueryRow(queryCheckID, id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar a existência do ID", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "ID não encontrado", http.StatusNotFound)
		return
	}

	// Continuar com a atualização no banco de dados
	queryUpdate := "UPDATE safisa.imposto SET uf = ?, icms = ?, ipi = ?, difal = ? WHERE id = ?"
	_, err = i.Db.Exec(queryUpdate, novoImposto.Uf, novoImposto.Icms, novoImposto.Ipi, novoImposto.Difal, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar dados no banco de dados", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Imposto atualizado com sucesso!")
}

func (i *Imposto) DeletarImposto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verificar se o ID já existe
	queryCheckID := "SELECT COUNT(*) FROM safisa.imposto WHERE id = ?"
	var count int
	err = i.Db.QueryRow(queryCheckID, id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar a existência do ID", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "ID não encontrado", http.StatusNotFound)
		return
	}

	// Continuar com a exclusão no banco de dados
	queryDelete := "DELETE FROM safisa.imposto WHERE id = ?"
	_, err = i.Db.Exec(queryDelete, id)
	if err != nil {
		http.Error(w, "Erro ao excluir dados do banco de dados", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Imposto deletado com sucesso!")
}

func (i *Imposto) ListarImpostoID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	query := "SELECT id, uf, icms, ipi, difal FROM safisa.imposto WHERE id = ?"
	row := i.Db.QueryRow(query, id)

	var imposto Imposto
	err = row.Scan(&imposto.Id, &imposto.Uf, &imposto.Icms, &imposto.Ipi, &imposto.Difal)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Imposto não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar dados do banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imposto)
}
