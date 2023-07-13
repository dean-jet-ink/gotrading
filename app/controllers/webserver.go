package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

type WebServer struct {
	address          string
	port             int
	templates        *template.Template
	candleController CandleController
}

func NewWebServer(address string, port int, dbConn *sql.DB) *WebServer {
	templateList := []string{
		"app/views/chart.html",
	}
	templates := template.Must(template.ParseFiles(templateList...))

	candleController := NewBitflyerCandleController(dbConn)

	return &WebServer{
		address:          address,
		port:             port,
		templates:        templates,
		candleController: candleController,
	}
}

func (ws *WebServer) viewHandler(w http.ResponseWriter, r *http.Request) {
	if err := ws.templates.ExecuteTemplate(w, "chart.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) Start() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc(ws.candleController.APIUrl(), ws.candleController.APIHandler)
	http.HandleFunc("/chart/", ws.viewHandler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", ws.address, ws.port), nil)
}
