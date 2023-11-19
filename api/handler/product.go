package handler

import (
	"encoding/json"
	"exam/models"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cast"
)

func (h *Handler) Product(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.CreateProduct(w, r)
	case "GET":
		var (
			id = r.URL.Query().Get("id")
		)
		if id != "" {
			h.ProductGetById(w, r)
		} else {
			h.ProductGetList(w, r)
		}
	case "DELETE":
		h.ProductDelete(w, r)
	case "PUT":
		h.ProductUpdate(w, r)
	default:
		fmt.Println("default")
	}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	var createProduct models.CreateProduct
	err := json.NewDecoder(r.Body).Decode(&createProduct)

	if err != nil {
		h.handleResponse(w, 500, "error while unmarsh request: "+err.Error())
		return
	}

	resp, err := h.storage.Product().Create(&createProduct)
	if err != nil {
		fmt.Println(createProduct)
		h.handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.handleResponse(w, http.StatusCreated, resp)
}

func (h *Handler) ProductGetById(w http.ResponseWriter, r *http.Request) {
	var prId = r.URL.Query().Get("id")
	resp, err := h.storage.Product().GetByID(&models.ProductPrimaryKey{Id: prId})
	if err != nil {
		h.handleResponse(w, 500, "Product does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) ProductGetList(w http.ResponseWriter, r *http.Request) {
	var offsetstr = r.URL.Query().Get("offset")
	var limitstr = r.URL.Query().Get("limit")
	var search = r.URL.Query().Get("search")

	var offset = cast.ToInt(offsetstr)
	var limit = cast.ToInt(limitstr)

	resp, err := h.storage.Product().GetList(&models.GetListProductRequest{Offset: offset, Limit: limit, Search: search})
	if err != nil {
		h.handleResponse(w, 500, "Product does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
func (h *Handler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	var ProductId = r.URL.Query().Get("id")

	resp, err := h.storage.Product().Delete(&models.ProductPrimaryKey{Id: ProductId})
	if err != nil {
		log.Println(err.Error())
		return
	}
	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	var product models.UpdateProduct

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resp, err := h.storage.Product().Update(&product)

	if err != nil {
		log.Println(err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
