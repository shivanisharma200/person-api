package main

import (
	httpPerson "person-api/internal/handler"
	servicePerson "person-api/internal/service/person"
	storePerson "person-api/internal/store/person"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	g := gofr.New()
	g.Server.ValidateHeaders = false

	store := storePerson.New()
	service := servicePerson.New(store)
	personHandler := httpPerson.New(service)

	g.GET("/persons/{id}", personHandler.GetByID)
	g.GET("/persons", personHandler.Get)
	g.POST("/persons", personHandler.Create)
	g.PUT("/persons/{id}", personHandler.Update)
	g.DELETE("/persons/{id}", personHandler.Delete)

	g.Start()
}
