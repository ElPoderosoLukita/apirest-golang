package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Usuario struct {
	// id       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      int    `json:"age"`
}

var id int = 0
var UsuariosMap map[string]*Usuario = make(map[string]*Usuario)

//GetUsers - GET - /api/users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var usuarios []Usuario

	for _, v := range UsuariosMap {
		usuarios = append(usuarios, *v)
	}

	objetoJson, err := json.Marshal(usuarios)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(objetoJson)

}

//PostUsers - POST - /api/users
func PostUsers(w http.ResponseWriter, r *http.Request) {
	var user Usuario

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	k := strconv.Itoa(id)
	UsuariosMap[k] = &user
	id++

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Usuario creado correctamente.")
}

//PutUsers - PUT - /api/users/{id}
func PutUsers(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario
	params := mux.Vars(r)
	k := params["id"]

	value, ok := UsuariosMap[k]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El id que nos ha solicitado no pertenece a ningún usuario.")
	}

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		panic(err)
	}

	value.Age = usuario.Age
	value.Lastname = usuario.Lastname
	value.Name = usuario.Name

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Usuario con id "+k+" actualizado correctamente.")
}

//DeleteUsers - DELETE - /api/users/{id}
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	k := params["id"]

	_, ok := UsuariosMap[k]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El id que nos ha solicitado no pertenece a ningún usuario.")
	}

	delete(UsuariosMap, k)
	fmt.Fprintf(w, "Usuario con id "+k+" eliminado correctamente.")

}

func main() {
	r := mux.NewRouter().StrictSlash(false)

	r.HandleFunc("/api/users", GetUsers).Methods("GET")
	r.HandleFunc("/api/users", PostUsers).Methods("POST")
	r.HandleFunc("/api/users/{id}", PutUsers).Methods("PUT")
	r.HandleFunc("/api/users/{id}", DeleteUsers).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}
	fmt.Println("Listening on http://localhost:8081...")
	server.ListenAndServe()
}
