package engine

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/adohong4/carZone/core"
	"github.com/adohong4/carZone/models"
	"github.com/adohong4/carZone/service"
	"github.com/adohong4/carZone/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		service: service,
	}
}

func (e *EngineHandler) GetEngineByID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "GetEngineByID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := e.service.GetEngineById(ctx, id)
	if err != nil {
		log.Printf("Error getting engine by ID: %v", err)
		if err.Error() == "Engine not found" {
			core.SendErrorResponse(w, core.NewNotFoundError("Engine not found").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Internal server error", utils.InternalServerError))
		return
	}

	core.NewOK("Engine retrieved successfully", resp).Send(w)
}

func (e *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "CreateEngine-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid request body").ErrorResponse)
		log.Println("Error reading request body: ", err)
		return
	}

	var engineReq models.EngineRequest
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		log.Println("Error unmarshalling request body: ", err)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid engine data").ErrorResponse)
		return
	}

	createdEngine, err := e.service.CreateEngine(ctx, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating engine: ", err)
		if err.Error() == "engine already exists" {
			core.SendErrorResponse(w, core.NewConflictRequestError("engine already exists").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Internal server error", utils.InternalServerError))
		return
	}

	core.NewCREATED("Engine created successfully", createdEngine).Send(w)
}

func (e *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateEngine-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid request body").ErrorResponse)
		return
	}

	var engineReq models.EngineRequest
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid engine data").ErrorResponse)
		return
	}

	updatedEngine, err := e.service.UpdateEngine(ctx, id, &engineReq)
	if err != nil {
		log.Printf("Error updating engine: %v", err)
		if err.Error() == "Engine not found" {
			core.SendErrorResponse(w, core.NewNotFoundError("Engine not found").ErrorResponse)
			return
		}
		core.SendErrorResponse(w, core.NewErrorResponse("Invalid server error", utils.InternalServerError))
		return
	}

	core.NewOK("Engine updated successfully", updatedEngine).Send(w)
}

func (e *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteEngine-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	deletedEngine, err := e.service.DeleteEngine(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting engine: ", err)
		response := map[string]string{"error": "Invalid ID or Engine not found"}
		jsonResponse, _ := json.Marshal(response)
		_, _ = w.Write(jsonResponse)
		return
	}

	// check if engine was deleted successfully
	if deletedEngine.EngineID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{"error": "Engine not found"}
		jsonResponse, _ := json.Marshal(response)
		_, _ = w.Write(jsonResponse)
		return
	}

	jsonResponse, err := json.Marshal(deletedEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response body: ", err)
		response := map[string]string{"error": "Internal server error"}
		jsonResponse, _ = json.Marshal(response)
		_, _ = w.Write(jsonResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}
