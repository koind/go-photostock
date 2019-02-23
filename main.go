package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
	"html/template"
	"fmt"
	"io/ioutil"
)

type ImageController struct {
}

func (imgC *ImageController) indexAction(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, nil)
}

func (imgC *ImageController) newAction(w http.ResponseWriter, r *http.Request) {
	file, info, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return
	}

	contentType := info.Header.Get("Content-Type")
	if !(contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif") {
		fmt.Printf("Wrong content type: %s", contentType)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	//imageName := fmt.Sprintf("%d.png", time.Now().Unix())

	ioutil.WriteFile(info.Filename, data, 0600)
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
	router.HandleFunc("/", imageController.indexAction).Methods("GET")
	router.HandleFunc("/images", imageController.newAction).Methods("POST")
	router.HandleFunc("/images/{id}", imageController.getOneByIdAction).Methods("GET")
	router.HandleFunc("/images/{id}/parts/{part_num}", imageController.getOneByIdAndPartIdAction).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
