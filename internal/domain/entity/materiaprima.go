package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// aço - inox - aluminio -bronze
// chumbaloy - aluminio 3320 - inox 316l - bronze3350
type MateriaPrima struct {
	Id         int
	Descrição  string
	Material   string
	Liga       string
	Kgmt       float64
	Valorkg    float64 // Massa kg por metro
	Imposto_id int
	Db         *sql.DB `json:"-"`
}

func NovaMateriaPrima(db *sql.DB) (*MateriaPrima, error) {
	return &MateriaPrima{
		Id:         0,
		Descrição:  "",
		Material:   "",
		Liga:       "",
		Kgmt:       0,
		Valorkg:    0,
		Imposto_id: 0,
		Db:         db,
	}, nil
}

//lebrar de fazer a função calculando por barras e por metro

func (mp *MateriaPrima) CriarMateriaPrima(w http.ResponseWriter, r *http.Request) {
	var materiaPrima MateriaPrima
	err := json.NewDecoder(r.Body).Decode(&materiaPrima)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Valida os campos obrigatórios
	if materiaPrima.Descrição == "" || materiaPrima.Material == "" || materiaPrima.Liga == "" || materiaPrima.Kgmt <= 0 || materiaPrima.Valorkg <= 0 || materiaPrima.Imposto_id == 0 {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	// Cria a consulta SQL para inserir a MateriaPrima
	query := `
		INSERT INTO safisa.materiaprima (id, descricao, material, liga, kgmt, valorkg, imposto_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Executa a consulta SQL
	_, err = mp.Db.Exec(query, materiaPrima.Id, materiaPrima.Descrição, materiaPrima.Material, materiaPrima.Liga, materiaPrima.Kgmt, materiaPrima.Valorkg, materiaPrima.Imposto_id)
	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Matéria-prima criada com sucesso!")

}

func (mp *MateriaPrima) ListarMateriaPrima(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, descricao, material, liga, kgmt, valorkg, imposto_id FROM safisa.materiaprima
	`

	rows, err := mp.Db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var materiasPrimas []MateriaPrima
	for rows.Next() {
		var materiaPrima MateriaPrima
		err := rows.Scan(&materiaPrima.Id, &materiaPrima.Descrição, &materiaPrima.Material, &materiaPrima.Liga, &materiaPrima.Kgmt, &materiaPrima.Valorkg, &materiaPrima.Imposto_id)
		if err != nil {
			http.Error(w, "Erro ao ler dados do banco de dados", http.StatusInternalServerError)
			return
		}
		materiasPrimas = append(materiasPrimas, materiaPrima)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(materiasPrimas)

}

func (mp *MateriaPrima) AtualizarMateriaPrima(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	var materiaPrima MateriaPrima
	err = json.NewDecoder(r.Body).Decode(&materiaPrima)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Valida os campos obrigatórios
	if materiaPrima.Descrição == "" || materiaPrima.Material == "" || materiaPrima.Liga == "" || materiaPrima.Kgmt <= 0 || materiaPrima.Valorkg <= 0 || materiaPrima.Imposto_id == 0 {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	// Verifica se o registro com o ID existe
	existsQuery := "SELECT COUNT(*) FROM safisa.materiaprima WHERE id = ?"
	var exists int
	err = mp.Db.QueryRow(existsQuery, id).Scan(&exists)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "Matéria-prima não encontrada", http.StatusNotFound)
		return
	}

	// Cria a consulta SQL para atualizar a MateriaPrima
	query := `
		UPDATE safisa.materiaprima
		SET descricao = ?, material = ?, liga = ?, kgmt = ?, valorkg = ?, imposto_id = ?
		WHERE id = ?
	`

	// Executa a consulta SQL
	_, err = mp.Db.Exec(query, materiaPrima.Descrição, materiaPrima.Material, materiaPrima.Liga, materiaPrima.Kgmt, materiaPrima.Valorkg, materiaPrima.Imposto_id, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Matéria-prima atualizada com sucesso!")
}
func (mp *MateriaPrima) DeletarMateriaPrima(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	// Verifica se o registro com o ID existe
	existsQuery := "SELECT COUNT(*) FROM safisa.materiaprima WHERE id = ?"
	var exists int
	err = mp.Db.QueryRow(existsQuery, id).Scan(&exists)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "Matéria-prima não encontrada", http.StatusNotFound)
		return
	}

	// Cria a consulta SQL para deletar a MateriaPrima
	query := "DELETE FROM safisa.materiaprima WHERE id = ?"

	// Executa a consulta SQL
	_, err = mp.Db.Exec(query, id)
	if err != nil {
		http.Error(w, "Erro ao deletar dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Matéria-prima deletada com sucesso!")
}

func (mp *MateriaPrima) ListarMateriaPrimaId(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	// Verifica se o registro com o ID existe
	existsQuery := "SELECT COUNT(*) FROM safisa.materiaprima WHERE id = ?"
	var exists int
	err = mp.Db.QueryRow(existsQuery, id).Scan(&exists)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "Matéria-prima não encontrada", http.StatusNotFound)
		return
	}

	// Cria a consulta SQL para buscar a MateriaPrima por ID
	query := `
	SELECT id, descricao, material, liga, kgmt, valorkg, imposto_id
	FROM safisa.materiaprima
	WHERE id = ?
`

	// Executa a consulta SQL
	var materiaPrima MateriaPrima
	err = mp.Db.QueryRow(query, id).Scan(&materiaPrima.Id, &materiaPrima.Descrição, &materiaPrima.Material, &materiaPrima.Liga, &materiaPrima.Kgmt, &materiaPrima.Valorkg, &materiaPrima.Imposto_id)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna a matéria-prima encontrada em formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(materiaPrima)
}
