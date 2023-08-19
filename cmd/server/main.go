package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"

	"github.com/wevyrtonantero/estoque/internal/domain/entity"
	
)

func main() {

	db, err := sql.Open("mysql", "root:wedeju180587@tcp(localhost:3306)/safisa")
	if err != nil {
		panic(err)

	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
		return
	}

	log.Println("Conexão com o banco de dados estabelecida")

	c, err := entity.NovaCategoria(db)
	if err != nil {
		panic(err)
	}
	cl, err := entity.NovoConsumivel(db)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Route("/categoria", func(r chi.Router) {
		r.Post("/criar", c.CriarCategoria)
		r.Get("/", c.ListarCategoria)
		r.Post("/atualizar", c.AtualizarCategoria)
		r.Delete("/deletar/{id}", c.DeletarCategoria)
		r.Get("/{id}", c.ListarCategoriaId)

	})

	r.Route("/consumivel", func(r chi.Router) {
		r.Post("/criar", cl.CriarConsumivel)
		r.Get("/", cl.ListarConsumivel)
		r.Post("/atualizar", cl.AtualizarConsumivel)
		r.Delete("/deletar{id}", cl.DeletarConsumivel)
		r.Get("/{id}", cl.ListarConsumivelId)

	})

	r.Route("/fornecedor", func(r chi.Router) {
		r.Put("/", entity.CriarFornecedor)
		r.Get("/", entity.ListarFornecedor)
		r.Put("/", entity.AtualizarFornecedor)
		r.Delete("/", entity.DeletarFornecedor)
		r.Get("/{id}", entity.ListarFornecedorId)

	})

	r.Route("/materiaprima", func(r chi.Router) {
		r.Put("/", entity.CriarMateriaPrima)
		r.Get("/", entity.ListarMateriaPrima)
		r.Put("/", entity.AtualizarMateriaPrima)
		r.Delete("/", entity.DeletarMateriaPrima)
		r.Get("/{id}", entity.ListarMateriaPrimaId)

	})

	r.Route("/peca", func(r chi.Router) {
		r.Put("/", entity.CriarPeca)
		r.Get("/", entity.ListarPeca)
		r.Put("/", entity.AtualizarPeca)
		r.Delete("/", entity.DeletarPeca)
		r.Get("/{id}", entity.ListarPecaId)

	})

	r.Route("/pedido", func(r chi.Router) {
		r.Put("/", entity.CriarPedido)
		r.Get("/", entity.ListarPedido)
		r.Put("/", entity.AtualizarPedido)
		r.Delete("/", entity.DeletarPedido)
		r.Get("/{id}", entity.ListarPedidoId)

	})

	r.Route("/submontagem", func(r chi.Router) {
		r.Put("/", entity.CriarSubMontagem)
		r.Get("/", entity.ListarSubMontagem)
		r.Put("/", entity.AtualizarSubMontagem)
		r.Delete("/", entity.DeletarSubMontagem)
		r.Get("/{id}", entity.ListarSubMontagemId)

	})

	println("Servidor rodando")

	http.ListenAndServe(":8080", r)

}
