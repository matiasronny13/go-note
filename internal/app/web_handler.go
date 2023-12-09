package app

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/matiasronny13/go-note/config"
	"github.com/matiasronny13/go-note/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type webHandler struct {
	AppConfig *config.AppConfig
	Service   AppService
}

type WebHandler interface {
	Default(ctx *gin.Context)
	GetById(ctx *gin.Context)
	DeleteById(c *gin.Context)
	Save(ctx *gin.Context)
	Create(ctx *gin.Context)
	Encrypt(ctx *gin.Context)
	Decrypt(ctx *gin.Context)
}

func NewWebHandler(appConfig *config.AppConfig, service AppService) *webHandler {
	return &webHandler{appConfig, service}
}

func (r *webHandler) BindState(c *gin.Context) *model.PageState {
	state := &model.PageState{}
	if err := c.Bind(state); err != nil {
		slog.Error("Error binding state", "message", err)
	} else {
		state.PageTitle = r.AppConfig.Web.Title
		state.PathId = c.Param("id")
		state.ShowDeleteButton = true
		state.IsEditMode = true
	}
	return state
}

func (r *webHandler) Default(c *gin.Context) {
	state := r.BindState(c)
	state.ShowDeleteButton = false
	c.HTML(http.StatusOK, "index", state)
}

func (r *webHandler) GetById(c *gin.Context) {
	state := r.BindState(c)
	r.Service.GetbyId(state)
	state.IsEditMode = false
	c.HTML(http.StatusOK, "index", state)
}

func (r *webHandler) DeleteById(c *gin.Context) {
	state := r.BindState(c)

	if state.Password == "" {
		state.Errors = append(state.Errors, errors.New("Password is required"))
	} else {
		r.Service.DeleteById(state)
		if len(state.Errors) == 0 {
			c.Header("HX-Redirect", "/")
		}
	}

	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Save(c *gin.Context) {
	state := r.BindState(c)

	if state.Password == "" {
		state.Errors = append(state.Errors, errors.New("Password is required"))
	} else {
		r.Service.Save(state)
		if len(state.Errors) == 0 {
			c.Header("HX-Redirect", "/"+state.Id)
		}
	}

	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Create(c *gin.Context) {
	state := r.BindState(c)

	if state.Password == "" {
		state.Errors = append(state.Errors, errors.New("Password is required"))
		state.ShowDeleteButton = false
	} else {
		r.Service.Create(state)
		if len(state.Errors) == 0 {
			c.Header("HX-Redirect", "/"+state.Id)
		}
	}

	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Encrypt(c *gin.Context) {
	state := r.BindState(c)

	if state.Password == "" {
		state.Errors = append(state.Errors, errors.New("Password is required"))
	} else {
		r.Service.EncryptMessage(state)
	}

	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Decrypt(c *gin.Context) {
	state := r.BindState(c)

	if state.Password == "" {
		state.Errors = append(state.Errors, errors.New("Password is required"))
	} else {
		r.Service.DecryptMessage(state)
	}

	c.HTML(http.StatusOK, "index_partial", state)
}
