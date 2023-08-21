package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Consumivel struct {
	Id           int
	Cod          string
	Descricao    string
	Categoria_id int
	Imposto_id   int
	Db           *sql.DB `json:"-"`
}

func NovoConsumivel(db *sql.DB) (*Consumivel, error) {
	return &Consumivel{
		Id:           0,
		Cod:          "",
		Descricao:    "",
		Categoria_id: 0,
		Imposto_id:   0,
		Db:           db,
	}, nil
}


func (cml *Consumivel) CriarConsumivel(w http.ResponseWriter, r *http.Request) {
	var consumivel Consumivel

	// Decodifica o JSON do corpo da solicitação para a estrutura Consumivel
	err := json.NewDecoder(r.Body).Decode(&consumivel)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Valida os campos obrigatórios
	if consumivel.Cod == "" || consumivel.Descricao == "" || consumivel.Categoria_id <= 0 || consumivel.Imposto_id <= 0 {
		http.Error(w, "Campos 'Cod', 'Descricao', 'Categoria_id' e 'Imposto_id' são obrigatórios e devem ser valores positivos", http.StatusBadRequest)
		return
	}

	// Cria a consulta SQL para inserir um novo Consumivel
	query := `
		INSERT INTO safisa.consumivel ( id, cod, descricao, categoria_id, imposto_id)
		VALUES (?, ?, ?, ?, ?)
	`

	// Executa a consulta SQL
	_, err = cml.Db.Exec(query, consumivel.Id, consumivel.Cod, consumivel.Descricao, consumivel.Categoria_id, consumivel.Imposto_id)
	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco de dados verifique se o Id é valido, e se foi digitado corretamente Impostos e Categoria", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Consumível criado com sucesso!")
}

func (cml *Consumivel) ListarConsumivel(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, cod, descricao, categoria_id, imposto_id FROM safisa.consumivel"
	rows, err := cml.Db.Query(query)
	if err != nil {
		http.Error(w, "Erro ao consultar dados no banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	consumiveis := []Consumivel{}

	for rows.Next() {
		var consumivel Consumivel
		err := rows.Scan(&consumivel.Id, &consumivel.Cod, &consumivel.Descricao, &consumivel.Categoria_id, &consumivel.Imposto_id)
		if err != nil {
			http.Error(w, "Erro ao ler dados do banco de dados", http.StatusInternalServerError)
			return
		}
		consumiveis = append(consumiveis, consumivel)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consumiveis)

}


func (cml *Consumivel) AtualizarConsumivel(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	// Consulta o banco de dados para verificar se o ID existe
	var count int
	err = cml.Db.QueryRow("SELECT COUNT(*) FROM safisa.consumivel WHERE id = ?", id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar o ID no banco de dados", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Consumível com ID especificado não encontrado", http.StatusNotFound)
		return
	}

	var consumivel Consumivel

	// Decodifica o JSON do corpo da solicitação para a estrutura Consumivel
	err = json.NewDecoder(r.Body).Decode(&consumivel)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Valida os demais campos obrigatórios
	if consumivel.Cod == "" || consumivel.Descricao == "" || consumivel.Categoria_id <= 0 || consumivel.Imposto_id <= 0 {
		http.Error(w, "Campos 'Cod', 'Descricao', 'Categoria_id' e 'Imposto_id' são obrigatórios e devem ser valores positivos", http.StatusBadRequest)
		return
	}

	// Cria a consulta SQL para atualizar o Consumivel
	query := `
		UPDATE safisa.consumivel
		SET cod = ?, descricao = ?, categoria_id = ?, imposto_id = ?
		WHERE id = ?
	`

	// Executa a consulta SQL
	_, err = cml.Db.Exec(query, consumivel.Cod, consumivel.Descricao, consumivel.Categoria_id, consumivel.Imposto_id, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Consumível atualizado com sucesso!")

}


func (cml *Consumivel) DeletarConsumivel(w http.ResponseWriter, r *http.Request)  {
	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	// Consulta o banco de dados para verificar se o ID existe
	var count int
	err = cml.Db.QueryRow("SELECT COUNT(*) FROM safisa.consumivel WHERE id = ?", id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar o ID no banco de dados", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Consumível com ID especificado não encontrado", http.StatusNotFound)
		return
	}

	// Cria a consulta SQL para excluir o Consumivel por ID
	query := "DELETE FROM safisa.consumivel WHERE id = ?"
	_, err = cml.Db.Exec(query, id)
	if err != nil {
		http.Error(w, "Erro ao excluir dados no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Consumível excluído com sucesso!")
}


func (cml *Consumivel) ListarConsumivelId(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido na URL", http.StatusBadRequest)
		return
	}

	// Consulta o banco de dados para buscar o Consumivel por ID
	query := "SELECT id, cod, descricao, categoria_id, imposto_id FROM safisa.consumivel WHERE id = ?"
	row := cml.Db.QueryRow(query, id)

	var consumivel Consumivel
	err = row.Scan(&consumivel.Id, &consumivel.Cod, &consumivel.Descricao, &consumivel.Categoria_id, &consumivel.Imposto_id)
	if err != nil {
		http.Error(w, "Erro ao buscar o Consumível no banco de dados", http.StatusInternalServerError)
		return
	}

	// Retorna o Consumivel encontrado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consumivel)
}
