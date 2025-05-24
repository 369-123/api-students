package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/369-123/api-students/schemas"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// GetStudents godoc
// @Summary      Get all students
// @Description  Retrieve all student records
// @Tags         students
// @Accept       json
// @Produce      json
// @Success      200 {array} schemas.StudentResponse
// @Failure      404 {string} string "not found"
// @Router       /students [get]
func (api *API) getStudents(c echo.Context) error {
	students, err := api.DB.GetStudents()
	if err != nil {
		return c.String(http.StatusNotFound, "Failed to get students")
	}

	return c.JSON(http.StatusOK, schemas.NewResponse(students))
}

// CreateStudent godoc
// @Summary      Create a new student
// @Description  Create a new student record
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        student body StudentRequest true "Student info"
// @Success      200 {string} string "create student"
// @Failure      400 {string} string "bad request"
// @Failure      500 {string} string "internal error"
// @Router       /students [post]
func (api *API) createStudents(c echo.Context) error {
	studentReq := StudentRequest{}
	if err := c.Bind(&studentReq); err != nil {
		return err
	}

	if err := studentReq.Validate(); err != nil {
		log.Error().Err(err).Msgf("[api] error validating struct")
		return c.String(http.StatusBadRequest, "Error validating student")
	}

	student := schemas.Student{
		Name:   studentReq.Name,
		Email:  studentReq.Email,
		CPF:    studentReq.CPF,
		Age:    studentReq.Age,
		Active: *studentReq.Active,
	}

	if err := api.DB.AddStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Error to create student")
	}
	return c.String(http.StatusOK, "create student")
}

// GetStudent godoc
// @Summary      Get student by ID
// @Description  Get a single student record by ID
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        id path int true "Student ID"
// @Success      200 {object} schemas.StudentResponse
// @Failure      404 {string} string "student not found"
// @Failure      500 {string} string "internal error"
// @Router       /students/{id} [get]
func (api *API) getStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get student ID")
	}

	student, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get student")
	}

	return c.JSON(http.StatusOK, student)
}

// UpdateStudent godoc
// @Summary      Update a student
// @Description  Update a student's details by ID
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        id path int true "Student ID"
// @Param        student body schemas.Student true "Student data to update"
// @Success      200 {object} schemas.Student
// @Failure      404 {string} string "student not found"
// @Failure      500 {string} string "internal error"
// @Router       /students/{id} [put]
func (api *API) updateStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get student ID")
	}

	receivedStudent := schemas.Student{}
	if err := c.Bind(&receivedStudent); err != nil {
		return err
	}

	updatingStudent, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get student")
	}

	student := updateStudentInfo(receivedStudent, updatingStudent)

	if err := api.DB.UpdateStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "failed to save student")
	}

	return c.JSON(http.StatusOK, student)
}

// DeleteStudent godoc
// @Summary      Delete a student
// @Description  Delete a student by ID
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        id path int true "Student ID"
// @Success      200 {object} schemas.Student
// @Failure      404 {string} string "student not found"
// @Failure      500 {string} string "internal error"
// @Router       /students/{id} [delete]
func (api *API) deleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get student ID")
	}

	student, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get student")
	}

	if err := api.DB.DeleteStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "failed to delete student")
	}

	return c.JSON(http.StatusOK, student)
}

func updateStudentInfo(receivedStudent, student schemas.Student) schemas.Student {
	if receivedStudent.Name != "" {
		student.Name = receivedStudent.Name
	}

	if receivedStudent.CPF != "" {
		student.CPF = receivedStudent.CPF
	}

	if receivedStudent.Email != "" {
		student.Email = receivedStudent.Email
	}

	if receivedStudent.Age > 0 {
		student.Age = receivedStudent.Age
	}

	if receivedStudent.Active != student.Active {
		student.Active = receivedStudent.Active
	}

	return student
}
