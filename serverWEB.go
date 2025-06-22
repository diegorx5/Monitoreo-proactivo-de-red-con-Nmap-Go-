package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Historico struct {
	ID      int
	IP      string
	Fecha   string
	Puertos string
	OS      string
	Conex   string
}

func main() {
	http.HandleFunc("/", handler)
	log.Println(" Servidor web en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./inventario.db")
	if err != nil {
		http.Error(w, "Error al abrir la base de datos", 500)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, ip, fecha, puertos, os, conex FROM historico ORDER BY id DESC`)
	if err != nil {
		http.Error(w, "Error en la consulta", 500)
		return
	}
	defer rows.Close()

	var historicos []Historico
	for rows.Next() {
		var h Historico
		rows.Scan(&h.ID, &h.IP, &h.Fecha, &h.Puertos, &h.OS, &h.Conex)
		historicos = append(historicos, h)
	}

	// Plantilla HTML simple
	const tpl = `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Dashboard NAC</title>
		<style>
			body { font-family: Arial; background: #f0f0f0; padding: 20px; }
			table { border-collapse: collapse; width: 100%; background: white; }
			th, td { border: 1px solid #ddd; padding: 8px; }
			th { background-color: #333; color: white; }
		</style>
	</head>
	<body>
		<h1>Histórico de Dispositivos Escaneados</h1>
		<table>
			<tr>
				<th>ID</th>
				<th>IP</th>
				<th>Fecha</th>
				<th>Puertos</th>
				<th>OS</th>
				<th>Conexiones</th>
			</tr>
			{{range .}}
			<tr>
				<td>{{.ID}}</td>
				<td>{{.IP}}</td>
				<td>{{.Fecha}}</td>
				<td>{{.Puertos}}</td>
				<td>{{.OS}}</td>
				<td><pre>{{.Conex}}</pre></td>
			</tr>
			{{end}}
		</table>
	</body>
	</html>`

	tmpl := template.Must(template.New("dashboard").Parse(tpl))
	tmpl.Execute(w, historicos)
}
