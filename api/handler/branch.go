package handler

import (
	"encoding/json"
	"exam/models"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cast"
)

func (h *Handler) Branch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.CreateBranch(w, r)
	case "GET":
		var (
			id = r.URL.Query().Get("id")
		)
		if id != "" {
			h.BranchGetById(w, r)
		} else {
			h.BranchGetList(w, r)
		}
	case "DELETE":
		h.BranchDelete(w, r)
	case "PUT":
		h.BranchUpdate(w, r)
	default:
		fmt.Println("default")
	}
}

func (h *Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	var createBranch models.CreateBranch
	err := json.NewDecoder(r.Body).Decode(&createBranch)

	if err != nil {
		h.handleResponse(w, 500, "error while unmarsh request: "+err.Error())
		return
	}

	resp, err := h.storage.Branch().Create(&createBranch)
	if err != nil {
		fmt.Println(createBranch)
		h.handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.handleResponse(w, http.StatusCreated, resp)
}

func (h *Handler) BranchGetById(w http.ResponseWriter, r *http.Request) {
	var brId = r.URL.Query().Get("id")
	fmt.Println("hello")
	resp, err := h.storage.Branch().GetByID(&models.BranchPrimaryKey{ID: brId})
	if err != nil {
		h.handleResponse(w, 500, "Branch does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) BranchGetList(w http.ResponseWriter, r *http.Request) {
	var offsetstr = r.URL.Query().Get("offset")
	var limitstr = r.URL.Query().Get("limit")
	var search = r.URL.Query().Get("search")

	name := r.URL.Query().Get("name")
	fromDate := r.URL.Query().Get("from_date")
	toDate := r.URL.Query().Get("to_date")

	var offset = cast.ToInt(offsetstr)
	var limit = cast.ToInt(limitstr)

	resp, err := h.storage.Branch().GetList(&models.GetListBranchRequest{Offset: offset, Limit: limit, Search: search, FromDate: fromDate, ToDate: toDate, Name: name})
	if err != nil {
		h.handleResponse(w, 500, "Branch does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
func (h *Handler) BranchDelete(w http.ResponseWriter, r *http.Request) {
	var branchId = r.URL.Query().Get("id")

	resp, err := h.storage.Branch().Delete(&models.BranchPrimaryKey{ID: branchId})
	if err != nil {
		log.Println(err.Error())
		return
	}
	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) BranchUpdate(w http.ResponseWriter, r *http.Request) {
	var branch models.UpdateBranch

	err := json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resp, err := h.storage.Branch().Update(&branch)

	if err != nil {
		log.Println(err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
