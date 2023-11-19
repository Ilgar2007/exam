package handler

import (
	"encoding/json"
	"exam/models"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cast"
)

func (h *Handler) Category(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.CreateCategory(w, r)
	case "GET":
		var (
			id = r.URL.Query().Get("id")
		)
		if id != "" {
			fmt.Println(id, "WORKING")
			h.CategoryGetById(w, r)
		} else {
			h.CategoryGetList(w, r)
		}
	case "DELETE":
		h.CategoryDelete(w, r)
	case "PUT":
		h.CategoryUpdate(w, r)
	default:
		fmt.Println("default")
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	var createCategory models.CreateCategory
	err := json.NewDecoder(r.Body).Decode(&createCategory)

	if err != nil {
		h.handleResponse(w, 500, "error while unmarsh request: "+err.Error())
		return
	}

	resp, err := h.storage.Category().Create(&createCategory)
	if err != nil {
		fmt.Println(createCategory)
		h.handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.handleResponse(w, http.StatusCreated, resp)
}

func (h *Handler) CategoryGetById(w http.ResponseWriter, r *http.Request) {
	var catId = r.URL.Query().Get("id")
	fmt.Println("hello")
	resp, err := h.storage.Category().GetById(&models.CategoryPrimaryKey{Id: catId})
	if err != nil {
		h.handleResponse(w, 500, "Category does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) CategoryGetList(w http.ResponseWriter, r *http.Request) {
	var offsetstr = r.URL.Query().Get("offset")
	var limitstr = r.URL.Query().Get("limit")
	var search = r.URL.Query().Get("search")

	var offset = cast.ToInt(offsetstr)
	var limit = cast.ToInt(limitstr)

	resp, err := h.storage.Category().GetList(&models.GetListCategoryRequest{Offset: offset, Limit: limit, Search: search})
	if err != nil {
		h.handleResponse(w, 500, "Category does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
func (h *Handler) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	var categoryId = r.URL.Query().Get("id")

	resp, err := h.storage.Category().Delete(&models.CategoryPrimaryKey{Id: categoryId})
	if err != nil {
		log.Println(err.Error())
		return
	}
	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	var category models.UpdateCategory

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resp, err := h.storage.Category().Update(&category)

	if err != nil {
		log.Println(err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
