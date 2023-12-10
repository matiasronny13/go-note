package app

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/matiasronny13/go-note/config"
	"github.com/matiasronny13/go-note/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type webHandler struct {
	AppConfig *config.AppConfig
	Service   AppService
}

type WebHandler interface {
	UnexpectedError(c *gin.Context, err any)
	AuthenticateUser() gin.HandlerFunc
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
		state.ShowDeleteButton = (state.PathId != "")
		state.IsEditMode = true
	}
	return state
}

func (r *webHandler) UnexpectedError(c *gin.Context, err any) {
	slog.Error("Unhandled exception", "error", err)
	state := c.MustGet("state").(*model.PageState)
	state.Errors = append(state.Errors, errors.New("Internal error, please contact admin."))
	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var state *model.PageState
		if !strings.HasPrefix(c.Request.RequestURI, "/assets") && c.Request.RequestURI != "/favicon.ico" {
			state = r.BindState(c)
			c.Set("state", state)
		}

		if c.Request.Method != "GET" && c.Request.Header["Content-Type"][0] == "application/x-www-form-urlencoded" {
			if state.Password == "" {
				state.Errors = append(state.Errors, errors.New("Password is required"))
			} else if state.PathId != "" {
				if passwordFromDb, err := r.Service.ValidatePassword(state.PathId, state.Password); err != nil {
					state.Errors = append(state.Errors, err)
				} else {
					state.Password = passwordFromDb // hashed
				}
			}

			if len(state.Errors) > 0 {
				c.HTML(http.StatusOK, "index_partial", state)
				c.Abort()
			}
		}

		c.Next()
	}
}

func (r *webHandler) Default(c *gin.Context) {
	c.HTML(http.StatusOK, "index", c.MustGet("state").(*model.PageState))
}

func (r *webHandler) GetById(c *gin.Context) {
	state := c.MustGet("state").(*model.PageState)
	r.Service.GetbyId(state)
	state.IsEditMode = false
	c.HTML(http.StatusOK, "index", state)
}

func (r *webHandler) DeleteById(c *gin.Context) {
	state := c.MustGet("state").(*model.PageState)

	r.Service.DeleteById(state)
	if len(state.Errors) == 0 {
		c.Header("HX-Redirect", "/")
	}

	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Save(c *gin.Context) {
	state := c.MustGet("state").(*model.PageState)

	r.Service.Save(state)
	if len(state.Errors) == 0 {
		c.Header("HX-Redirect", "/"+state.Id)
	}

	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Create(c *gin.Context) {
	state := c.MustGet("state").(*model.PageState)

	r.Service.Create(state)
	if len(state.Errors) == 0 {
		c.Header("HX-Redirect", "/"+state.Id)
	}

	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Encrypt(c *gin.Context) {
	state := c.MustGet("state").(*model.PageState)
	r.Service.EncryptMessage(state)
	c.HTML(http.StatusOK, "index_partial", state)
}

func (r *webHandler) Decrypt(c *gin.Context) {
	state := c.MustGet("state").(*model.PageState)
	r.Service.DecryptMessage(state)
	c.HTML(http.StatusOK, "index_partial", state)
}
