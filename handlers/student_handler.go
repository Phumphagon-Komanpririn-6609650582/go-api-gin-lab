package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/student-api/models"
	"example.com/student-api/repositories"
	"example.com/student-api/services"
)

type StudentHandler struct {
	Service *services.StudentService
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
	students, err := h.Service.GetStudents()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve students")
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id := c.Param("id")
	student, err := h.Service.GetStudentByID(id)
	if err != nil {

		if err == repositories.ErrNotFound {
			respondWithError(c, http.StatusNotFound, "Student not found")
			return
		}

		respondWithError(c, http.StatusInternalServerError, "Failed to get student")
		return
	}
	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid input: id, name must not be empty and gpa must be between 0.00 and 4.00")
		return
	}

	if err := h.Service.CreateStudent(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) UpdateStudentHandler(c *gin.Context) {
	id := c.Param("id")
	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid input data: check required fields and GPA range (0.00-4.00)")
		return
	}

	err := h.Service.UpdateStudent(id, student)
	if err != nil {
		if err == repositories.ErrNotFound {
			respondWithError(c, http.StatusNotFound, "Student not found")
			return
		}
		respondWithError(c, http.StatusInternalServerError, "Failed to update student")
		return
	}

	student.Id = id
	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) DeleteStudentHandler(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.DeleteStudent(id)
	if err != nil {
		if err == repositories.ErrNotFound {
			respondWithError(c, http.StatusNotFound, "Student not found")
			return
		}
		respondWithError(c, http.StatusInternalServerError, "Failed to delete student")
		return
	}

	c.Status(http.StatusNoContent)
}
