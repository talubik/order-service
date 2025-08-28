package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"html/template"
	"log"
	"myapp/src/internal/memory"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
)


type OrderHandler struct{
	mem *memory.Memory

}

func NewOrderHandler(m *memory.Memory) *OrderHandler{
	return &OrderHandler{
		mem : m,
	}
}

func (h *OrderHandler) HandleOrder(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id:= vars["id"]
	order ,err := h.mem.Get(id)
	if err!= nil{
		http.Error(w,"Order was not found", http.StatusBadRequest)
		return
	}
	isCurl := strings.Contains(r.UserAgent(), "curl") || 
	r.Header.Get("Accept") == "application/json"
	if isCurl {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
	return
	}
	abspath  := filepath.Join(getProjectRoot(),"template","order.html")
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		log.Printf("Template not found at: %s", abspath)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	tmpl , err := template.ParseFiles(abspath)
	if err!=nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
	}
	w.Header().Set("Content_Type","text/html")
	if err := tmpl.Execute(w, order); err!= nil{
		log.Printf("Template error: %v", err)
	}
	//json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) ShowOrderSearchForm(w http.ResponseWriter, r *http.Request) {
	abspath := filepath.Join(getProjectRoot(),"template","order_search.html")
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		log.Printf("Template not found at: %s", abspath)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
    tmpl := template.Must(template.ParseFiles(abspath))
    tmpl.Execute(w, nil)
}

func getProjectRoot() string{
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	return projectRoot
}

