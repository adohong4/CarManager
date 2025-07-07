package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/adohong4/carZone/models"
	"github.com/adohong4/carZone/service"
	"github.com/gorilla/mux"
)

type CarHandler struct {
	service service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{
		service: service,
	}
}

func (h *CarHandler) GetCarById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.service.GetCarById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response: ", err)
	}
}

func (h *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"

	resp, err := h.service.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response: ", err)
	}
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Err : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var carReq models.CarRequest
	err = json.Unmarshal(body, &carReq)
	if err != nil {
		log.Println("Error unmarshalling request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createdCar, err := h.service.CreateCar(ctx, &carReq)
	if err != nil {
		log.Println("Error creating car: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(createdCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response body: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// write the response body
	_, _ = w.Write(responseBody)
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var carReq models.CarRequest
	err = json.Unmarshal(body, &carReq)
	if err != nil {
		log.Println("Error unmarshalling request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updatedCar, err := h.service.UpdateCar(ctx, id, &carReq)
	if err != nil {
		log.Println("Error updating car: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := json.Marshal(updatedCar)
	if err != nil {
		log.Println("Error marshalling response body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(resBody)
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]

	deletedCar, err := h.service.DeleteCar(ctx, id)
	if err != nil {
		log.Println("Error deleting car: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := json.Marshal(deletedCar)
	if err != nil {
		log.Println("Error marshalling response body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(resBody)
	if err != nil {
		log.Println("Error writing response: ", err)
	}
}
