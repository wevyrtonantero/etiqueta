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

type Emproducao struct {
	Id           int
	Tempomaquina int
	Cnc          string
	Operador     string
	Peca_id      int
	Db           *sql.DB `json:"-"`
}



func NovoEmProducao(db *sql.DB) (*Emproducao, error) {
	return &Emproducao{
		Id:           0,
		Tempomaquina: 10,
		Cnc:          "",
		Operador:     "",
		Peca_id:      0,
		Db:           db,
	}, nil

}

func (ep *Emproducao) CriarEmProducao(w http.ResponseWriter, r *http.Request) {
	var novaProducao Emproducao
	err := json.NewDecoder(r.Body).Decode(&novaProducao)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
		log.Printf("Erro ao ler o corpo da solicitação: %v", err)
		return
	}
	defer r.Body.Close()

	// Verifique se o ID da peça é válido
	if novaProducao.Id <= 0 {
		http.Error(w, "ID da peça inválido", http.StatusBadRequest)
		log.Printf("ID da peça inválido: %d", novaProducao.Id)
		return
	}

	// Verifique se o ID da peça existe no banco de dados
	query := "SELECT COUNT(*) FROM safisa.emproducao WHERE id = ?"
	var count int
	err = ep.Db.QueryRow(query, novaProducao.Id).Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao verificar o ID da peça", http.StatusInternalServerError)
		log.Printf("Erro ao verificar o ID da peça: %v", err)
		return
	}
	if count > 0 {
		http.Error(w, "ID existente ", http.StatusNotFound)
		log.Printf("ID da peça ja existe: %d", novaProducao.Peca_id)
		return
	}

	// Agora podemos inserir a produção no banco de dados
	insertQuery := "INSERT INTO safisa.emproducao (id, tempomaquina, peca_id, cnc, operador) VALUES (?, ?, ?, ?, ?)"
	_, err = ep.Db.Exec(insertQuery, novaProducao.Id, novaProducao.Tempomaquina, novaProducao.Peca_id, novaProducao.Cnc, novaProducao.Operador)
	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco de dados", http.StatusInternalServerError)
		log.Printf("Erro ao inserir dados no banco de dados: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Produção criada com sucesso!")
	log.Printf("Produção criada com sucesso: %+v", novaProducao)

}

func (ep *Emproducao) ListarEmProducao(w http.ResponseWriter, r *http.Request) {
			rows, err := ep.Db.Query("SELECT * FROM safisa.emproducao")
		if err != nil {
			http.Error(w, "Erro ao ler dados do banco de dados", http.StatusInternalServerError)
			log.Printf("Erro ao ler dados do banco de dados: %v", err)
			return
		}
		defer rows.Close()
	
		var producoes []Emproducao
		for rows.Next() {
			var producao Emproducao
			var tempomaquina sql.NullInt64
			var cnc sql.NullString
			var operador sql.NullString
			err := rows.Scan(&producao.Id, &tempomaquina, &cnc, &operador, &producao.Peca_id)
			if err != nil {
				http.Error(w, "Erro ao ler dados do banco de dados", http.StatusInternalServerError)
				log.Printf("Erro ao ler dados do banco de dados: %v", err)
				return
			}
			if tempomaquina.Valid {
				producao.Tempomaquina = int(tempomaquina.Int64)
			}
			if cnc.Valid {
				producao.Cnc = cnc.String
			}
			if operador.Valid {
				producao.Operador = operador.String
			}
			producoes = append(producoes, producao)
		}
	
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(producoes)
		if err != nil {
			http.Error(w, "Erro ao escrever resposta JSON", http.StatusInternalServerError)
			log.Printf("Erro ao escrever resposta JSON: %v", err)
			return
		}

}

func (ep *Emproducao) AtualizarEmProducao( w http.ResponseWriter, r *http.Request){
// Extrai o ID do parâmetro da URL
idStr := chi.URLParam(r, "id")
id, err := strconv.Atoi(idStr)
if err != nil {
	http.Error(w, "ID inválido", http.StatusBadRequest)
	log.Printf("ID inválido: %v", err)
	return
}

// Verifica se o ID da produção existe no banco de dados
query := "SELECT COUNT(*) FROM safisa.emproducao WHERE id = ?"
var count int
err = ep.Db.QueryRow(query, id).Scan(&count)
if err != nil {
	http.Error(w, "Erro ao verificar o ID da produção", http.StatusInternalServerError)
	log.Printf("Erro ao verificar o ID da produção: %v", err)
	return
}
if count == 0 {
	http.Error(w, "ID da produção não encontrado", http.StatusNotFound)
	log.Printf("ID da produção não encontrado: %d", id)
	return
}

var atualizacao Emproducao
err = json.NewDecoder(r.Body).Decode(&atualizacao)
if err != nil {
	http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusBadRequest)
	log.Printf("Erro ao ler o corpo da solicitação: %v", err)
	return
}
defer r.Body.Close()

// Atualiza os campos no banco de dados, incluindo o campo peca_id
updateQuery := "UPDATE safisa.emproducao SET tempomaquina = ?, cnc = ?, operador = ?, peca_id = ? WHERE id = ?"
_, err = ep.Db.Exec(updateQuery, atualizacao.Tempomaquina, atualizacao.Cnc, atualizacao.Operador, atualizacao.Peca_id, id)
if err != nil {
	http.Error(w, "Erro ao atualizar dados no banco de dados, Verifique de Id da peça é valido", http.StatusInternalServerError)
	log.Printf("Erro ao atualizar dados no banco de dados: %v", err)
	return
}

w.WriteHeader(http.StatusOK)
fmt.Fprintln(w, "Produção atualizada com sucesso!")
log.Printf("Produção atualizada com sucesso: %d", id)

}



func (ep *Emproducao) ListarEmProducaoId(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do parâmetro da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		log.Printf("ID inválido: %s", idStr)
		return
	}

	// Consulta o banco de dados para obter os dados da produção
	query := "SELECT id, tempomaquina, peca_id, cnc, operador FROM safisa.emproducao WHERE id = ?"
	var producao Emproducao
	err = ep.Db.QueryRow(query, id).Scan(&producao.Id, &producao.Tempomaquina, &producao.Peca_id, &producao.Cnc, &producao.Operador)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Produção não encontrada", http.StatusNotFound)
			log.Printf("Produção não encontrada para o ID: %d", id)
			return
		}
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		log.Printf("Erro ao consultar o banco de dados: %v", err)
		return
	}

	// Retorna os dados da produção como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producao)
}

func (ep *Emproducao) DeletarEmProducao(w http.ResponseWriter, r *http.Request){

// Extrair o ID do parâmetro da URL
idStr := chi.URLParam(r, "id")
id, err := strconv.Atoi(idStr)
if err != nil {
	http.Error(w, "ID inválido", http.StatusBadRequest)
	log.Printf("ID inválido: %s", idStr)
	return
}

// Verificar se o ID existe no banco de dados
query := "SELECT COUNT(*) FROM safisa.emproducao WHERE id = ?"
var count int
err = ep.Db.QueryRow(query, id).Scan(&count)
if err != nil {
	http.Error(w, "Erro ao verificar o ID da produção", http.StatusInternalServerError)
	log.Printf("Erro ao verificar o ID da produção: %v", err)
	return
}
if count == 0 {
	http.Error(w, "ID da produção não encontrado", http.StatusNotFound)
	log.Printf("ID da produção não encontrado: %d", id)
	return
}

// Deletar o registro do banco de dados
deleteQuery := "DELETE FROM safisa.emproducao WHERE id = ?"
_, err = ep.Db.Exec(deleteQuery, id)
if err != nil {
	http.Error(w, "Erro ao deletar a produção do banco de dados", http.StatusInternalServerError)
	log.Printf("Erro ao deletar a produção do banco de dados: %v", err)
	return
}

w.WriteHeader(http.StatusOK)
fmt.Fprintf(w, "Produção deletada com sucesso!")
log.Printf("Produção com ID %d foi deletada", id)


}
