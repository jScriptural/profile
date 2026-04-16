package api

import (
	"encoding/json"
	"log"
	"net/http"
	"profile/internal/service"
	"profile/internal/models"
	"io"
	"errors"
	"strings"
	"unicode"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		svc: s,
	}

}

func (h *Handler) HandleProfileCreation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data models.PostData
	


	if err := h.decode(r.Body, &data); err != nil {
		h.sendResponse(
			w,
			http.StatusUnprocessableEntity,
			"error",
			"Invalid type",
			nil,
			nil,
		)
		return
	}


	log.Printf("postData: %v",data)
	if strings.TrimSpace(data.Name) == "" {
		h.sendResponse(
			w,
			http.StatusBadRequest,
			"error",
			"Missing or empty name",
			nil,
			nil,
		);
		return
	}
	data.Name = strings.ToLower(strings.TrimSpace(data.Name))

	p, isNew, err := h.svc.GetOrCreateProfile(r.Context(), data.Name)

	if err != nil {
		if errors.Is(err, models.Err502) {
			h.sendResponse(
				w,
				http.StatusBadGateway,
				"error",
				"Genderize returned an invalid response",
				nil,
				nil,
			);
			return
		}

		h.sendResponse(
			w,
			http.StatusInternalServerError,
			"error","Upstream or server error",
			nil,
			nil,
		);
		return
	}

	if isNew {
		h.sendResponse(
			w,
			http.StatusCreated,
			"success",
			"",
			nil,
			p,
		);
		return
	}

	status := "success"
	msg := "Profile already exists"
	h.sendResponse(
		w,
		http.StatusOK,
		status,
		msg,
		nil,
		p,
	)
	return
}



func  (h *Handler) HandleProfileRetrievalByID(w http.ResponseWriter, r *http.Request){
	id := h.removeAllWhitespaces(r.PathValue("uuid"));
	
	if id == "" {
		h.sendResponse(
			w,
			http.StatusBadRequest,
			"error",
			"Missing or empty id",
			nil,
			nil,
		)
		return;
	}

	p,err := h.svc.RetrieveProfileByID(r.Context(),id);
	if err != nil {
		if errors.Is(err,models.ErrNoRows){
			h.sendResponse(
				w,
				http.StatusNotFound,
				"error",
				"Profile not found",
				nil,
				nil,
			)
			return
		}

		h.sendResponse(
			w,
			http.StatusInternalServerError,
			"error",
			"Unexpected server error",
			nil,
			nil,
		)

		return;
	}

	h.sendResponse(w,http.StatusOK,"success","",nil,p)

}


func (h *Handler) HandleAllProfileRetrievalWithFilter(w http.ResponseWriter, r *http.Request){
	values := r.URL.Query();

	log.Printf("values: %#v\n",values)

	gender := values.Get("gender")
	countryID := values.Get("country_id")
	ageGroup := values.Get("age_group");



	var count int = 23;
	h.sendResponse(
		w,
		http.StatusOK,
		"success",
		"testing",
		&count,
		struct{gender string;ageGroup string;countryID string}{gender,ageGroup,countryID},
	)

}




/****************************************  
*                                       *
*            HELPER FUNCS               *
*                                       *  
*****************************************/

func (h *Handler) decode(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

func (h *Handler) encode(w io.Writer, data any) error {
	log.Printf("%v",data)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}


func (h *Handler) sendResponse(w http.ResponseWriter, statusCode int, status, msg string, count *int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := models.APIResponse{
		Status:  status,
		Message: msg,
		Data:    data,
		Count: count,
	}


	if err := h.encode(w, resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

func (h *Handler)removeAllWhitespaces(s string) string {
	trim := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)

	return trim;
}
