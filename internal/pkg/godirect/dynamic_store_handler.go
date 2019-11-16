package godirect

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type CreateRequest struct {
	Code int    `json:"code"`
	Url  string `json:"url"`
}

type CreateResponse struct {
	Id        string `json:"id"`
	Code      int    `json:"code"`
	TargetUrl string `json:"target_url"`
	SourceUrl string `json:"source_url"`
}

func DynamicDirectHandlerFunc(directorURL *url.URL, store *DynamicDirectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		cRequest := &CreateRequest{}
		err := json.NewDecoder(r.Body).Decode(cRequest)
		if err != nil {
			http.Error(w, "Bad Request\n"+err.Error(), http.StatusBadRequest)
			return
		}
		targetUrl, err := url.Parse(cRequest.Url)
		if err != nil {
			http.Error(w, "Bad Request\n"+err.Error(), http.StatusBadRequest)
			return
		}

		if !in(targetUrl.Scheme, "http", "https") {
			http.Error(w, "Bad Request\nInvalid url scheme (http,https)", http.StatusBadRequest)
			return
		}

		if !in(cRequest.Code, http.StatusTemporaryRedirect, http.StatusMovedPermanently) {
			http.Error(w, "Bad Request\nInvalid code (301,307)", http.StatusBadRequest)
			return
		}

		direct, err := store.CreateAndAdd(cRequest.Code, targetUrl)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		sourceUrl, err := directorURL.Parse(direct.Path())
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		response := CreateResponse{
			Id:        direct.Path(),
			Code:      direct.Code(),
			TargetUrl: direct.URL(),
			SourceUrl: sourceUrl.String(),
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), http.StatusInternalServerError)
		}

	}
}
