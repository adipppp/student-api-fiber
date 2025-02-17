package controllers

import (
	"studentapifiber/constants"
	"studentapifiber/db"
	"studentapifiber/dto"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type StudentHandler struct {
}

func (*StudentHandler) GetStudents(c *fiber.Ctx) error {
	pool, err := db.GetDbPool()
	if err != nil {
		log.Errorf("error getting database pool: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.InternalServerError,
			Message: "Internal Server Error",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	rows, err := pool.Query(c.Context(), "select npm, name from student")
	if err != nil {
		log.Errorf("error querying database: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.InternalServerError,
			Message: "Internal Server Error",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	students := make([]dto.Student, 0)
	for rows.Next() {
		var npm, name string
		err = rows.Scan(&npm, &name)
		if err != nil {
			log.Errorf("error scanning row: %v", err)
			errorResponse := dto.ErrorResponse{
				Code:    constants.InternalServerError,
				Message: "Internal Server Error",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
		}
		student := dto.Student{NPM: npm, Name: name}
		students = append(students, student)
	}

	return c.Status(fiber.StatusOK).JSON(students)
}

func (*StudentHandler) GetStudentById(c *fiber.Ctx) error {
	npmParam := c.Params("npm")
	if npmParam == "" {
		errorResponse := dto.ErrorResponse{
			Code:    constants.InvalidNpm,
			Message: "NPM is required",
		}
		return c.Status(fiber.StatusNotFound).JSON(errorResponse)
	}

	pool, err := db.GetDbPool()
	if err != nil {
		log.Errorf("error getting database pool: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.InternalServerError,
			Message: "Internal Server Error",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	var npm, name string
	err = pool.QueryRow(c.Context(), "select npm, name from student where npm=$1", npmParam).Scan(&npm, &name)
	if err != nil {
		log.Errorf("error querying database: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.NotFound,
			Message: "Student data not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(errorResponse)
	}

	student := dto.Student{NPM: npm, Name: name}
	return c.Status(fiber.StatusOK).JSON(student)
}

func (*StudentHandler) PostStudent(c *fiber.Ctx) error {
	student := new(dto.Student)
	err := c.BodyParser(student)
	if err != nil || student.NPM == "" || student.Name == "" {
		log.Errorf("error parsing request body: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.InternalServerError,
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	pool, err := db.GetDbPool()
	if err != nil {
		log.Errorf("error getting database pool: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.InternalServerError,
			Message: "Internal Server Error",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	nowTimestamp := time.Now().Format(time.RFC3339)
	_, err = pool.Exec(c.Context(), "insert into student(npm, name, created_at, modified_at) values($1, $2, $3, $4)", student.NPM, student.Name, nowTimestamp, nowTimestamp)
	if err != nil {
		log.Errorf("error inserting student data: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.StudentAlreadyExists,
			Message: "Student data already exists",
		}
		return c.Status(fiber.StatusConflict).JSON(errorResponse)
	}

	return c.Status(fiber.StatusCreated).JSON(student)
}

func (*StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	npmParam := c.Params("npm")
	if npmParam == "" {
		errorResponse := dto.ErrorResponse{
			Code:    constants.InvalidNpm,
			Message: "NPM is required",
		}
		return c.Status(fiber.StatusNotFound).JSON(errorResponse)
	}

	pool, err := db.GetDbPool()
	if err != nil {
		log.Errorf("error getting database pool: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.InternalServerError,
			Message: "Internal Server Error",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	_, err = pool.Exec(c.Context(), "delete from student where npm=$1", npmParam)
	if err != nil {
		log.Errorf("error deleting student data: %v", err)
		errorResponse := dto.ErrorResponse{
			Code:    constants.InternalServerError,
			Message: "Internal Server Error",
		}
		return c.Status(fiber.StatusNotFound).JSON(errorResponse)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
