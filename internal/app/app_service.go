package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/matiasronny13/go-note/config"
	"github.com/matiasronny13/go-note/internal/pkg/model"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/crypto/bcrypt"
)

type appService struct {
	AppConfig    *config.AppConfig
	DbClient     *supabase.Client
	CryptoClient CryptoService
}

type AppService interface {
	GetbyId(result *model.PageState)
	DeleteById(data *model.PageState)
	ValidatePassword(id string, inputPassword string) (initialPassword string, err error)

	Save(data *model.PageState)
	Create(data *model.PageState)
	EncryptMessage(data *model.PageState)
	DecryptMessage(data *model.PageState)
}

func NewAppService(appConfig *config.AppConfig, dbClient *supabase.Client, crypto CryptoService) *appService {
	return &appService{appConfig, dbClient, crypto}
}

func GenerateRandomId(length int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}

	return string(result), nil
}

func (r *appService) ValidatePassword(id string, inputPassword string) (initialPassword string, err error) {
	var resultSet *model.PageState
	r.DbClient.From("notes").Select("password", "", false).Eq("id", id).Single().ExecuteTo(&resultSet)
	if !CheckPasswordHash(inputPassword, resultSet.Password) {
		return "", errors.New("Invalid Password")
	}
	return resultSet.Password, nil
}

func (r *appService) GetbyId(result *model.PageState) {
	var queryResult []model.PageState

	if _, err := r.DbClient.From("notes").Select("id, content, is_encrypted", "", false).Eq("id", result.PathId).ExecuteTo(&queryResult); err != nil {
		slog.Error("Get by Id", "message", err)
	}

	if len(queryResult) > 0 {
		result.Content = queryResult[0].Content
		result.IsEncrypted = queryResult[0].IsEncrypted
	} else {
		result.ShowDeleteButton = false
		result.Errors = append(result.Errors, errors.New("Record not found"))
	}
}

func (r *appService) DeleteById(state *model.PageState) {
	r.DbClient.From("notes").Delete("", "").Eq("id", state.PathId).Execute()
}

func (r *appService) GenerateNewId() (result string, err error) {
	for i := 0; i < 3; i++ {
		if result, err = GenerateRandomId(5); err == nil {
			a, _, _ := r.DbClient.From("notes").Select("id", "", false).Eq("id", result).Single().ExecuteString()
			if len(a) == 0 {
				return
			}
		}
	}

	return "", errors.New("Failed to generate new id")
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *appService) CheckDuplicateId(id string) error {
	var found []*model.PageState
	r.DbClient.From("notes").Select("id", "", false).Eq("id", id).ExecuteTo(&found)
	if len(found) > 0 {
		return errors.New(fmt.Sprintf("Id \"%s\" has been used by existing record", id))
	}
	return nil
}

func (r *appService) Save(state *model.PageState) {
	if state.Id == "" {
		state.Id = state.PathId
	} else if state.PathId != state.Id {
		if err := r.CheckDuplicateId(state.Id); err != nil {
			state.Errors = append(state.Errors, err)
			return
		}
		go func() {
			r.DbClient.From("notes").Delete("", "").Eq("id", state.PathId).Execute()
		}()
	}

	r.DbClient.From("notes").Upsert(&state, "", "", "").Execute()
}

func (r *appService) Create(state *model.PageState) {
	if state.Id == "" {
		state.Id, _ = GenerateRandomId(5)
	} else if err := r.CheckDuplicateId(state.Id); err != nil {
		state.Errors = append(state.Errors, err)
		return
	}

	state.Password = HashPassword(state.Password)

	if _, _, err := r.DbClient.From("notes").Insert(&state, false, "", "", "").Execute(); err != nil {
		state.Errors = append(state.Errors, err)
	}
}

func (r *appService) EncryptMessage(state *model.PageState) {
	var err error
	if state.Content, err = r.CryptoClient.Encrypt(state.Content, state.Password); err != nil {
		state.Errors = append(state.Errors, err)
	} else {
		state.IsEncrypted = true
	}
}

func (r *appService) DecryptMessage(state *model.PageState) {
	var err error
	if state.Content, err = r.CryptoClient.Decrypt(state.Content, state.Password); err != nil {
		state.Errors = append(state.Errors, err)
	} else {
		state.IsEncrypted = false
	}
}
