package entity

// cotovelo - conexão de ar - conexão de óleo - empurradores.

import (
	"database/sql"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Categoria struct {
	Id        int
	Categoria string
	Peca_id    []int
	consumivel_id []int
	Db        *sql.DB `json:"-"`
}

func NovaCategoria(db *sql.DB) (*Categoria, error) {
	return &Categoria{
		Id:        0,
		Categoria: "",
		Db:        db,
	}, nil
}

func (c *Categoria) CriarCategoria(w http.ResponseWriter, r *http.Request) {

	var categoria Categoria
	err := json.NewDecoder(r.Body).Decode(&categoria)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	if categoria.Categoria == "" {
		http.Error(w, "O campo 'categoria' é obrigatório", http.StatusBadRequest)
		return
	}

	var count int
	err = c.Db.QueryRow("SELECT COUNT(*) FROM safisa.Categoria WHERE id = ?", categoria.Id).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		http.Error(w, "ID já existe no banco de dados", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO safisa.Categoria (id, categoria) VALUES (?, ?)"
	_, err = c.Db.Exec(query, categoria.Id, categoria.Categoria)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dados inseridos com sucesso!")

}

func (c *Categoria) ListarCategoria(w http.ResponseWriter, r *http.Request) {

	query := "SELECT id, categoria FROM safisa.Categoria"
	rows, err := c.Db.Query(query)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var categorias []Categoria
	for rows.Next() {
		var categoria Categoria
		err := rows.Scan(&categoria.Id, &categoria.Categoria)
		if err != nil {
			log.Fatal(err)
			return
		}
		categorias = append(categorias, categoria)
	}

	json.NewEncoder(w).Encode(categorias)
}

func (c *Categoria) AtualizarCategoria(w http.ResponseWriter, r *http.Request) {
	var categoria Categoria
	err := json.NewDecoder(r.Body).Decode(&categoria)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verifica se a categoria existe antes de atualizar
	var categoriaBanco int
	err = c.Db.QueryRow("SELECT COUNT(*) FROM safisa.Categoria WHERE id = ?", categoria.Id).Scan(&categoriaBanco)
	if err != nil {
		log.Fatal(err)
	}

	if categoriaBanco == 0 {
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}

	// Atualiza a categoriaDb
	query := "UPDATE safisa.Categoria SET categoria = ? WHERE id = ?"
	_, err = c.Db.Exec(query, categoria.Categoria, categoria.Id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode("Nome atualizada com sucesso!")

}

func (c *Categoria) DeletarCategoria(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM safisa.Categoria WHERE id = ?"
	_, err = c.Db.Exec(query, idInt)
	if err != nil {
		http.Error(w, "Erro ao deletar a categoria", http.StatusInternalServerError)
		return
	}

	fmt.Println("Categoria deletada com sucesso!")
}

func (c *Categoria) ListarCategoriaId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("bateu na rota GET Listar por id a CategoriaDb")

	id := chi.URLParam(r, "id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id inválido")
		return
	}

	if idint <= 0 {

		fmt.Fprintf(w, "id deve ser maior que ZERO")
		return

	}

	stmt, err := c.Db.Prepare("SELECT * FROM safisa.Categoria WHERE id =?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer stmt.Close()

	var categoria Categoria
	err = stmt.QueryRow(idint).Scan(&categoria.Id, &categoria.Categoria)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if categoria.Id == 0 {
		json.NewEncoder(w).Encode("CategoriaDb nao Existe")
		return
	}

	json.NewEncoder(w).Encode(categoria)
}
