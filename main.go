package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
)

type ImageController struct {
}

func (imgC *ImageController) newAction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func (imgC *ImageController) getOneByIdAction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func (imgC *ImageController) getOneByIdAndPartIdAction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func main() {
	imageController := ImageController{}

	router := mux.NewRouter()
	router.HandleFunc("/images", imageController.newAction).Methods("POST")
	router.HandleFunc("/images/{id}", imageController.getOneByIdAction).Methods("GET")
	router.HandleFunc("/images/{id}/parts/{part_num}", imageController.getOneByIdAndPartIdAction).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
