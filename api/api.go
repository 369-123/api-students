package api

import (
	"github.com/369-123/api-students/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	Echo *echo.Echo
	DB   *db.StudentHandler
}

func NewServer() *API {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dataBase := db.Init()
	studentDB := db.NewStudentHandler(dataBase)

	return &API{
		Echo: e,
		DB:   studentDB,
	}
}

func (api *API)ConfigureRoutes() {
	//routes
	api.Echo.GET("/students", api.getStudents)
	api.Echo.POST("/students", api.createStudents)
	api.Echo.GET("/students/:id", api.getStudent)
	api.Echo.PUT("/students/:id", api.updateStudent)
	api.Echo.DELETE("/students/:id", api.deleteStudent)
}

func (api *API) Start() error {
	// Start server
 return api.Echo.Start(":8080")

}