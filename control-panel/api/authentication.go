package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	api "platformlab/controlpanel/api/dto"
	"platformlab/controlpanel/service"

	"gorm.io/gorm"
)

type Authentication struct {
	userService service.User
}

func (a *Authentication) Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input api.LoginRequestDto
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			fmt.Println(err.Error())
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		user, err := a.userService.VerifyAuthenticationCredentials(&input.Email, &input.Password)
		if err != nil {
			log.Println("[AuthenticationAPI] user invalid credentials: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "credentials.invaild"})
			return
		}

		response := api.NewLoginSuccessResponseDto(user)
		json.NewEncoder(w).Encode(&response)
	}
}

func AuthenticationRESTApi(db *gorm.DB) *Authentication {
	return &Authentication{userService: service.User{Db: db}}
}
