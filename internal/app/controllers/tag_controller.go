package controllers

import (
	"net/http"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/services"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/logger"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/utils"
	pbTag "github.com/ponyjackal/go-microservice-boilerplate/proto/tag"

	"github.com/bufbuild/protovalidate-go"
	"github.com/gin-gonic/gin"
)

type TagController struct {
	tagService *services.TagService
	validator  *protovalidate.Validator
}

func NewTagController(tagService *services.TagService) *TagController {
	// Create a new validator
	validator, err := protovalidate.New()
	if err != nil {
		logger.Errorf("failed to initialize validator: %v", err)
	}

	return &TagController{
		tagService: tagService,
		validator:  validator,
	}
}

// GetTags godoc
// @Summary Retrieve a list of tags
// @Description Get a list of tags filtered by the name parameter
// @Tags Tags
// @Accept json
// @Produce json
// @Param name query string false "Name of the tag to filter by"
// @Success 200 {object} pbTag.GetTagsResponse "Successful retrieval of tags"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags [get]
func (c *TagController) GetTags(ctx *gin.Context) {
	name := ctx.Query("name")

	response, err := c.tagService.GetTags(name)
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetTagById godoc
// @Summary Retrieve a tag
// @Description Get tag by id from the database
// @Tags Tags
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} pbTag.Tag "Successfully retrieved a tag"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags/{id} [get]
func (c *TagController) GetTagById(ctx *gin.Context) {
	id := ctx.Param("id")

	response, err := c.tagService.GetTagById(&pbTag.TagId{
		Id: id,
	})
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// SaveTag godoc
// @Summary Add a new tag
// @Description Add a tag with the provided information
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body pbTag.SaveTagRequest true "Tag Object"
// @Success 201 {object} models.Tag "Successfully created tag"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags [post]
func (c *TagController) SaveTag(ctx *gin.Context) {
	var tagReq pbTag.SaveTagRequest

	if err := ctx.BindJSON(&tagReq); err != nil {
		logger.Errorf("Invalid request: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err})
		return
	}
	// Validate the request
	if err := c.validator.Validate(&tagReq); err != nil {
		logger.Errorf("Invalid request: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err})
		return
	}

	response, err := c.tagService.SaveTag(&tagReq)
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// UpdateTag godoc
// @Summary Update tag
// @Description Update a tag with the provided information
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Param tag body pbTag.SaveTagRequest true "Tag Object"
// @Success 200 {object} pbTag.Tag "Successfully updated a tag"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags/{id} [put]
func (c *TagController) UpdateTag(ctx *gin.Context) {
	id := ctx.Param("id")

	var tagReq pbTag.SaveTagRequest
	if err := ctx.BindJSON(&tagReq); err != nil {
		logger.Errorf("Invalid request: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err})
		return
	}

	tag, err := c.tagService.UpdateTag(&pbTag.UpdateTagRequest{
		Id:     id,
		TagReq: &tagReq,
	})
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &tag)
}

// DeleteTag godoc
// @Summary Delete tag
// @Description Delete a tag with the provided information
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags/{id} [delete]
func (c *TagController) DeleteTag(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.tagService.DeleteTag(&pbTag.TagId{
		Id: id,
	})
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}
