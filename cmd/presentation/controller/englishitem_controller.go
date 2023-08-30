package controller

import (
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnglishItemController interface {
	Create(c *gin.Context)
	Proposal(c *gin.Context)
}

type EnglishItemGinController struct {
	proposalUse usecase.ProposalEnglishItemUsecase
	createUse   usecase.CreateEnglishItemUsecase
}

func NewEnglishItemController(proposalUse usecase.ProposalEnglishItemUsecase, createUse usecase.CreateEnglishItemUsecase) EnglishItemController {
	return &EnglishItemGinController{
		proposalUse: proposalUse,
		createUse:   createUse,
	}
}

func (ec *EnglishItemGinController) Proposal(c *gin.Context) {
	req := &dto.ProposalEnglishItemRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := ec.proposalUse.Proposal(req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) Create(c *gin.Context) {
	req := &dto.CreateEnglishItemRequest{}
	if err := c.BindJSON(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userId, err := userId(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	req.UserId = userId
	resp, err := ec.createUse.Create(req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}
