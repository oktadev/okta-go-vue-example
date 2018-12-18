package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/klebervirgilio/vue-crud-app-with-golang/pkg/core"
	"github.com/klebervirgilio/vue-crud-app-with-golang/pkg/kudo"
)

type Service struct {
	repo   core.Repository
	Router http.Handler
}

func New(repo core.Repository) Service {
	service := Service{
		repo: repo,
	}

	router := httprouter.New()
	router.GET("/kudos", service.Index)
	router.POST("/kudos", service.Create)
	router.DELETE("/kudos/:id", service.Delete)
	router.PUT("/kudos/:id", service.Update)

	service.Router = UseMiddlewares(router)

	return service
}

func (s Service) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	service := kudo.NewService(s.repo, r.Context().Value("userId").(string))
	kudos, err := service.GetKudos()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kudos)
}

func (s Service) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	service := kudo.NewService(s.repo, r.Context().Value("userId").(string))
	payload, _ := ioutil.ReadAll(r.Body)

	githubRepo := kudo.GitHubRepo{}
	json.Unmarshal(payload, &githubRepo)

	kudo, err := service.CreateKudoFor(githubRepo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(kudo)
}

func (s Service) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	service := kudo.NewService(s.repo, r.Context().Value("userId").(string))

	repoID, _ := strconv.Atoi(params.ByName("id"))
	githubRepo := kudo.GitHubRepo{RepoID: int64(repoID)}

	_, err := service.RemoveKudo(githubRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s Service) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	service := kudo.NewService(s.repo, r.Context().Value("userId").(string))
	payload, _ := ioutil.ReadAll(r.Body)

	githubRepo := kudo.GitHubRepo{}
	json.Unmarshal(payload, &githubRepo)

	kudo, err := service.UpdateKudoWith(githubRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kudo)
}
