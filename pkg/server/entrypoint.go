package server

import (
	"fmt"
	"github.com/Datariah/psychic-guacamole/internal"
	"github.com/Datariah/psychic-guacamole/internal/secrets"
	"github.com/Datariah/psychic-guacamole/pkg/server/protected"
	"github.com/Datariah/psychic-guacamole/pkg/server/public"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	negronilog "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func getRouter() *mux.Router {
	r := public.GetRouter()
	protectedRouter := protected.GetRouter()

	secretValues, err := secrets.GetSecretValues("psychic-guacamole", internal.AwsRegion)

	if err != nil {
		log.Panicf("PANIC! %v", err)
	}

	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte((*secretValues)["JWT_SECRET"].(string)), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(protectedRouter))
	r.PathPrefix("/api/v1").Handler(an)

	return r
}

func Run() {
	router := getRouter()

	n := negroni.Classic()

	n.Use(negronilog.NewMiddleware())
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%d", 8081))
}
