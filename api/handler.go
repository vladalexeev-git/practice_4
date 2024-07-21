package api

import (
	"city-server/models"
	"city-server/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type handler struct {
	CityRepo repository.CityRepo
	router   *mux.Router
}

func NewHandler(r *mux.Router, cr repository.CityRepo) *handler {
	return &handler{
		CityRepo: cr,
		router:   r,
	}
}

func (h *handler) GetCity(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "city id is not passed", http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	city, err := h.CityRepo.GetCityById(cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(city)
	if err != nil {
		log.Println("failed to encode city", err)
	}
}

func (h *handler) ChangeCity(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "city id is not passed", http.StatusBadRequest)
		return
	}
	cid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("failed to convert city id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	n := query.Get("name")
	p := query.Get("population")

	var city models.City

	city.ID = cid
	city.Name = n
	city.Population, err = strconv.Atoi(p)
	if err != nil {
		log.Println("failed to convert population", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.CityRepo.UpdateCity(city)
	if err != nil {
		// todo: проверить что есть такой город
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) AddCity(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	n := query.Get("name")
	p := query.Get("population")

	if n == "" || p == "" {
		http.Error(w, "name or population is not passed", http.StatusBadRequest)
		return
	}

	ps, err := strconv.Atoi(p)
	if err != nil {
		log.Println("failed to convert population", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var city models.City
	city.Name = n
	city.Population = ps

	cityId, err := h.CityRepo.CreateCity(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(cityId)
	if err != nil {
		log.Println("failed to encode city", err)
	}
}

func (h *handler) GetAllCities(w http.ResponseWriter, r *http.Request) {
	cities, err := h.CityRepo.GetAllCities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(cities)
	if err != nil {
		log.Println("failed to encode cities", err)
	}
}
