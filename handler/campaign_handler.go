package handler

import (
	"fmt"
	"net/http"
	"rmzstartup/helper"
	model "rmzstartup/model/entity"
	"rmzstartup/service"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service service.CampaignService
}

func NewCampaignHandler(service service.CampaignService) *campaignHandler {
	return &campaignHandler{service: service}
}

// api/v1/campaings

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID := c.Query("user_id")

	// FindALL()
	if userID == "" {
		campaigns, err := h.service.GetCampaigns("")
		if err != nil {
			response := helper.APIResponse("Error get campaigns", http.StatusBadRequest, "error", helper.FormatCampaigns(campaigns))
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("List of campaign", http.StatusOK, "success", helper.FormatCampaigns(campaigns))
		c.JSON(http.StatusOK, response)
		return
	}

	// FindByID()
	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error get campaigns", http.StatusBadRequest, "error", helper.FormatCampaigns(campaigns))
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaign", http.StatusOK, "success", helper.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input helper.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", helper.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input helper.CreateCampaignInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(model.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully create campaign", http.StatusCreated, "success", helper.FormatCampaign(newCampaign))
	c.JSON(http.StatusCreated, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID helper.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData helper.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(model.User)

	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully update campaign", http.StatusOK, "success", helper.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImageCampaign(c *gin.Context) {
	var input helper.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(model.User)
	input.User = currentUser
	userID := currentUser.Id

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("failed to upload campaign image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// currentUser = c.MustGet("currentUser").(model.User)
	// userID := currentUser.Id.String()
	path := fmt.Sprintf("campaign-images/%s-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("failed to upload campaign image", http.StatusInternalServerError, "error", data)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
		"user_id":     userID,
	}
	response := helper.APIResponse("campaign image successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
