package controller

import (
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnglishItemController interface {
	Proposal(c *gin.Context)
}

type EnglishItemGinController struct {
	pu usecase.ProposalEnglishItemUsecase
}

func NewEnglishItemController(pu usecase.ProposalEnglishItemUsecase) EnglishItemController {
	return &EnglishItemGinController{
		pu: pu,
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

	resp, err := ec.pu.Proposal(req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
