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

type Categoria struct {
	Id        int
	Categoria string
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
		log.Println("Erro ao decodificar o JSON:", err)
		return
	}

	if categoria.Categoria == "" {
		http.Error(w, "O campo 'categoria' é obrigatório", http.StatusBadRequest)
		log.Println("Campo 'categoria' é obrigatório")
		return
	}

	query := "INSERT INTO safisa.categoria (id, categoria) VALUES (?, ?)"
	_, err = c.Db.Exec(query, categoria.Id, categoria.Categoria)
	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco de dados", http.StatusInternalServerError)
		log.Println("Erro ao inserir dados no banco de dados:", err)
		return
	}

}

func (c *Categoria) ListarCategoria(w http.ResponseWriter, r *http.Request) {
	fmt.Println("listar categoria ok")

	query := "SELECT id, categoria FROM safisa.categoria"
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
	err = c.Db.QueryRow("SELECT COUNT(*) FROM safisa.categoria WHERE id = ?", categoria.Id).Scan(&categoriaBanco)
	if err != nil {
		log.Fatal(err)
	}

	if categoriaBanco == 0 {
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}

	// Atualiza a categoriaDb
	query := "UPDATE safisa.categoria SET categoria = ? WHERE id = ?"
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

	query := "DELETE FROM safisa.categoria WHERE id = ?"
	_, err = c.Db.Exec(query, idInt)
	if err != nil {
		http.Error(w, "Erro ao deletar a categoria", http.StatusInternalServerError)
		return
	}

	fmt.Println("Categoria deletada com sucesso!")
}

func (c *Categoria) ListarCategoriaId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("bateu na rota GET Listar por id a Categoria")

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

	stmt, err := c.Db.Prepare("SELECT * FROM safisa.categoria WHERE id =?")
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
