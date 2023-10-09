package controller

import (
	"english/cmd/presentation/errhandle"
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnglishItemController interface {
	Create(c *gin.Context)
	Proposal(c *gin.Context)
	GetByUserIdAndContent(c *gin.Context)
}

type EnglishItemGinController struct {
	proposalUse usecase.ProposalEnglishItemUsecase
	createUse   usecase.CreateEnglishItemUsecase
	getUse      usecase.GetEnglishItemUsecase
}

func NewEnglishItemController(proposalUse usecase.ProposalEnglishItemUsecase, createUse usecase.CreateEnglishItemUsecase, getUse usecase.GetEnglishItemUsecase) EnglishItemController {
	return &EnglishItemGinController{
		proposalUse: proposalUse,
		createUse:   createUse,
		getUse:      getUse,
	}
}

func (ec *EnglishItemGinController) Proposal(c *gin.Context) {
	req := &dto.ProposalEnglishItemRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	resp, err := ec.proposalUse.Proposal(req)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) Create(c *gin.Context) {
	req := &dto.CreateEnglishItemRequest{}
	if err := c.BindJSON(req); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	userId, err := userId(c)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	req.UserId = userId
	if err := ec.createUse.Create(req); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.Status(http.StatusCreated)
}

func (ec *EnglishItemGinController) GetByUserIdAndContent(c *gin.Context) {
	content := c.Query("content")

	userId, err := userId(c)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	resp, err := ec.getUse.GetByUserIdAndContent(userId, content)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}
