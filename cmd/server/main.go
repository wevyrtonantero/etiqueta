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
	cto, err := entity.NovoCusto(db)
	if err != nil {
		panic(err)
	}
	i, err := entity.NovoImposto(db)
	if err != nil {
		panic(err)

	}

	cml, err := entity.NovoConsumivel(db)
	if err != nil {
		panic(err)

	}

	f, err := entity.NovoFornecedor(db)
	if err != nil {
		panic(err)

	}

	mp, err := entity.NovaMateriaPrima(db)
	if err != nil {
		panic(err)

	}

p, err := entity.NovaPeca(db)
	if err != nil {
		panic(err)

	}

	ep, err := entity.NovoEmProducao(db)
	if err != nil {
		panic(err)

	}

	r := chi.NewRouter()
	r.Route("/categoria", func(r chi.Router) {
		r.Post("/criar", c.CriarCategoria)
		r.Get("/", c.ListarCategoria)
		r.Post("/atualizar/{id}", c.AtualizarCategoria)
		r.Delete("/deletar/{id}", c.DeletarCategoria)
		r.Get("/{id}", c.ListarCategoriaId)

	})

	r.Route("/consumivel", func(r chi.Router) {
		r.Post("/criar", cml.CriarConsumivel)
		r.Get("/", cml.ListarConsumivel)
		r.Post("/atualizar/{id}", cml.AtualizarConsumivel)
		r.Delete("/deletar/{id}", cml.DeletarConsumivel)
		r.Get("/{id}", cml.ListarConsumivelId)

	})

	r.Route("/custo", func(r chi.Router) {
		r.Post("/criar", cto.CriarCusto)
		r.Get("/", cto.ListarCusto)
		r.Post("/atualizar/{id}", cto.AtualizarCusto)
		r.Delete("/deletar/{id}", cto.DeletarCusto)
		r.Get("/{id}", cto.ListarCustoID)

	})

	r.Route("/emproducao", func(r chi.Router) {
		r.Post("/criar", ep.CriarEmProducao)
		r.Get("/", ep.ListarEmProducao)
		r.Post("/atualizar/{id}", ep.AtualizarEmProducao)
		r.Delete("/deletar/{id}", ep.DeletarEmProducao)
		r.Get("/{id}", ep.ListarEmProducaoId)

	})

	r.Route("/imposto", func(r chi.Router) {
		r.Post("/criar", i.CriarImposto)
		r.Get("/", i.ListarImposto)
		r.Post("/atualizar/{id}", i.AtualizarImposto)
		r.Delete("/deletar/{id}", i.DeletarImposto)
		r.Get("/{id}", i.ListarImpostoID)

	})

	r.Route("/fornecedor", func(r chi.Router) {
		r.Post("/criar", f.CriarFornecedor)
		r.Get("/", f.ListarFornecedor)
		r.Post("/atualizar/{id}", f.AtualizarFornecedor)
		r.Delete("/deletar/{id}", f.DeletarFornecedor)
		r.Get("/{id}", f.ListarFornecedorId)

	})

	r.Route("/materiaprima", func(r chi.Router) {
		r.Post("/criar", mp.CriarMateriaPrima)
		r.Get("/", mp.ListarMateriaPrima)
		r.Post("/atualizar/{id}", mp.AtualizarMateriaPrima)
		r.Delete("/deletar/{id}", mp.DeletarMateriaPrima)
		r.Get("/{id}", mp.ListarMateriaPrimaId)

	})

	r.Route("/peca", func(r chi.Router) {
		r.Post("/criar", p.CriarPeca)
		r.Get("/", p.ListarPeca)
		r.Post("/atualizar/{id}", p.AtualizarPeca)
		r.Delete("/deletar/{id}", p.DeletarPeca)
		r.Get("/{id}", p.ListarPecaId)

	})

	r.Route("/pedido", func(r chi.Router) {
		r.Post("/", entity.CriarPedido)
		r.Get("/", entity.ListarPedido)
		r.Post("/", entity.AtualizarPedido)
		r.Delete("/", entity.DeletarPedido)
		r.Get("/{id}", entity.ListarPedidoId)

	})

	r.Route("/submontagem", func(r chi.Router) {
		r.Post("/", entity.CriarSubMontagem)
		r.Get("/", entity.ListarSubMontagem)
		r.Post("/", entity.AtualizarSubMontagem)
		r.Delete("/", entity.DeletarSubMontagem)
		r.Get("/{id}", entity.ListarSubMontagemId)

	})

	println("Servidor rodando")

	http.ListenAndServe(":8080", r)

}
