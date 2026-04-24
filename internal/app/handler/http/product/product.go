package hproduct

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/AlexBond702/catalog-service/internal/app/entity"
	rhandler "github.com/AlexBond702/catalog-service/internal/app/handler/http"
	"github.com/AlexBond702/catalog-service/internal/app/service"
	"github.com/AlexBond702/catalog-service/internal/pkg/binding"
	"github.com/AlexBond702/catalog-service/internal/pkg/http/httph"
)

type handler struct {
	svcService service.Product
}

func NewHandler(svcService service.Product) rhandler.Product {
	return &handler{svcService: svcService}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestProductCreate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.SendError(w, r, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}
	category, err := h.svcService.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.SendError(w, r, http.StatusBadRequest, err)
		default:
			httph.SendError(w, r, http.StatusInternalServerError, err)
		}
		return
	}
	resp := entity.ResponseCategory{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
	httph.SendJSON(w, http.StatusCreated, resp)
}

func (h *handler) GetByGUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["guid"])
	if err != nil {
		httph.SendError(w, r, http.StatusBadRequest, err)
		return
	}

	product, err := h.svcService.GetByGUID(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.SendError(w, r, http.StatusNotFound, err)
			httph.ErrorApply(r, err)
			return
		}
		httph.SendJSON(w, http.StatusInternalServerError, err)
		return
	}
	resp := entity.ResponseProduct{
		GUID:         product.GUID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
	httph.SendJSON(w, http.StatusOK, resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["guid"])
	if err != nil {
		httph.SendError(w, r, http.StatusBadRequest, err)
		return
	}
	var req entity.RequestProductUpdate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.SendError(w, r, http.StatusBadRequest, entity.ErrIncorrectParameters)
	}
	product, err := h.svcService.Update(r.Context(), guid, req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.SendError(w, r, http.StatusNotFound, err)
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.SendError(w, r, http.StatusBadRequest, err)
		default:
			httph.SendError(w, r, http.StatusInternalServerError, err)
		}
		return
	}
	resp := entity.RequestProductUpdate{
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
	}
	httph.SendJSON(w, http.StatusOK, resp)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.Parse(vars["guid"])
	if err != nil {
		httph.SendError(w, r, http.StatusBadRequest, err)
	}
	if err := h.svcService.Delete(r.Context(), guid); err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.SendError(w, r, http.StatusNotFound, err)
		case errors.Is(err, entity.ErrCategoryHasProducts):
			httph.SendError(w, r, http.StatusBadRequest, err)
		default:
			httph.SendError(w, r, http.StatusInternalServerError, err)
		}
		return
	}
	httph.SendJSON(w, http.StatusNoContent, nil)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.svcService.List(r.Context())
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.SendError(w, r, http.StatusNotFound, err)
			return
		}
		httph.SendJSON(w, http.StatusInternalServerError, err)
		return
	}
	resp := make([]entity.Product, len(products))
	for _, prod := range products {
		resp = append(resp, entity.Product{
			GUID:         prod.GUID,
			Name:         prod.Name,
			Description:  prod.Description,
			Price:        prod.Price,
			CategoryGUID: prod.CategoryGUID,
			CreatedAt:    prod.CreatedAt,
			UpdatedAt:    prod.UpdatedAt,
		})
	}
	httph.SendJSON(w, http.StatusOK, resp)
}
