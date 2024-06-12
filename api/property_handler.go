package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"propmanager/internal/app/model"
	"propmanager/internal/app/service"
)

// PropertyHandler represents the property handler.
type PropertyHandler struct {
	propertyService *service.PropertyService
	s3Service       *service.S3Service
}

// NewPropertyHandler returns a new property handler.
func NewPropertyHandler(propertyService *service.PropertyService, s3Service *service.S3Service) *PropertyHandler {
	return &PropertyHandler{propertyService: propertyService, s3Service: s3Service}
}

// GetAllProperties godoc
// @Summary Get all properties
// @Description Get all properties
// @Tags Properties
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Property
// @Failure 500 {object} map[string]string
// @Router /properties [get]
func (h *PropertyHandler) GetAllProperties(c *gin.Context) {
	properties, err := h.propertyService.GetAllProperties()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, properties)
}

// GetProperty godoc
// @Summary Get a property
// @Description Get a property by ID
// @Tags Properties
// @Accept  json
// @Produce  json
// @Param id path int true "Property ID"
// @Success 200 {object} model.Property
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /properties/{id} [get]
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	property, err := h.propertyService.GetProperty(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, property)
}

// CreateProperty godoc
// @Summary Create a property
// @Description Create a new property
// @Tags Properties
// @Accept  json
// @Produce  json
// @Param property body model.Property true "Property"
// @Success 201 {object} model.Property
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /properties [post]
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	var property model.Property
	if err := c.BindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.propertyService.CreateProperty(&property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, property)
}

// UpdateProperty godoc
// @Summary Update a property
// @Description Update a property
// @Tags Properties
// @Accept  json
// @Produce  json
// @Param id path int true "Property ID"
// @Param property body model.Property true "Property"
// @Success 200 {object} model.Property
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /properties/{id} [put]
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var property model.Property
	if err := c.BindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	property.ID = uint(id)

	if err := h.propertyService.UpdateProperty(&property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, property)
}

// DeleteProperty godoc
// @Summary Delete a property
// @Description Delete a property by ID
// @Tags Properties
// @Accept  json
// @Produce  json
// @Param id path int true "Property ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /properties/{id} [delete]
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.propertyService.DeleteProperty(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// UploadImage godoc
// @Summary Upload an image
// @Description Upload an image for a property
// @Tags Properties
// @Accept  multipart/form-data
// @Produce  json
// @Param id path int true "Property ID"
// @Param file formData file true "Image file"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /properties/{id}/images [post]
func (h *PropertyHandler) UploadImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileBytes := make([]byte, file.Size)
	_, err = openedFile.Read(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileName := fmt.Sprintf("%d-%s", id, file.Filename)
	url, err := h.s3Service.UploadImage(c.Request.Context(), fileBytes, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	property, err := h.propertyService.GetProperty(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	property.Images = append(property.Images, model.Image{URL: url})

	if err := h.propertyService.UpdateProperty(&property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url": url})
}

// DeleteImage godoc
// @Summary Delete an image
// @Description Delete an image associated with a property
// @Tags Properties
// @Accept  json
// @Produce  json
// @Param id path int true "Property ID"
// @Param image_id path int true "Image ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /properties/{id}/images/{image_id} [delete]
func (h *PropertyHandler) DeleteImage(c *gin.Context) {
	propertyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.propertyService.DeleteImage(uint(propertyID), uint(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
