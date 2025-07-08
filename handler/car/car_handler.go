package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/adohong4/carZone/core"
	"github.com/adohong4/carZone/models"
	"github.com/adohong4/carZone/service"
	"github.com/adohong4/carZone/utils"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
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
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "GetCarById-Handler")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.service.GetCarById(ctx, id)
	if err != nil {
		log.Printf("Error getting car by ID: %v", err)
		if err.Error() == "car not found" {
			core.SendErrorResponse(w, core.NewNotFoundError("Car not found").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Internal server error", utils.InternalServerError))
		return
	}
	core.NewOK("Car retrieved successfully", resp).Send(w)
}

func (h *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "GetCarByBrand-Handler")
	defer span.End()

	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"

	resp, err := h.service.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		log.Printf("Error getting car by brand: %v", err)
		if err.Error() == "Car not found" {
			core.SendErrorResponse(w, core.NewNotFoundError("No cars found for the given brand").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Internal server error", utils.InternalServerError))
		return
	}

	core.NewOK("Cars retrieved successfully", resp).Send(w)
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "CreateCar-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid request body").ErrorResponse)
		return
	}

	var carReq models.CarRequest
	if err = json.Unmarshal(body, &carReq); err != nil {
		log.Printf("Error unmarshalling request body: ", err)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid car data").ErrorResponse)
		return
	}

	createdCar, err := h.service.CreateCar(ctx, &carReq)
	if err != nil {
		log.Println("Error creating car: ", err)
		if err.Error() == "Car already exists" {
			core.SendErrorResponse(w, core.NewConflictRequestError("Car already exists").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Internal server error", utils.InternalServerError))
		return
	}

	core.NewCREATED("Car created successfully", createdCar).Send(w)
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateCar-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid request body").ErrorResponse)
		return
	}

	var carReq models.CarRequest
	err = json.Unmarshal(body, &carReq)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid car data").ErrorResponse)
		return
	}

	updatedCar, err := h.service.UpdateCar(ctx, id, &carReq)
	if err != nil {
		log.Printf("Error updating car: %v", err)
		if err.Error() == "Car not found" {
			core.SendErrorResponse(w, core.NewNotFoundError("Car not found").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Internal server error", utils.InternalServerError))
		return
	}

	core.NewOK("Car updated successfully", updatedCar).Send(w)
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteCar-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	_, err := h.service.DeleteCar(ctx, id)
	if err != nil {
		log.Printf("Error deleting car: %v", err)
		if err.Error() == "Car not found" {
			core.SendErrorResponse(w, core.NewNotFoundError("Car not found").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Internal server error", utils.InternalServerError))
		return
	}

	core.NewOK("Car deleted successfully", nil).Send(w)
}
