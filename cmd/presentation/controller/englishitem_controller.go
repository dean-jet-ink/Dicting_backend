package controller

import (
	"english/cmd/presentation/errhandle"
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"english/myerror"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnglishItemController interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Proposal(c *gin.Context)
	ProposalTranslation(c *gin.Context)
	ProposalExplanation(c *gin.Context)
	ProposalExample(c *gin.Context)
	GetByUserId(c *gin.Context)
	GetById(c *gin.Context)
	GetRequiredExp(c *gin.Context)
}

type EnglishItemGinController struct {
	proposalUse       usecase.ProposalEnglishItemUsecase
	createUse         usecase.CreateEnglishItemUsecase
	getUse            usecase.GetEnglishItemUsecase
	updateUse         usecase.UpdateEnglishItemUsecase
	deleteUse         usecase.DeleteEnglishItemUsecase
	getRequiredExpUse usecase.GetRequiredExpUsecase
}

func NewEnglishItemController(proposalUse usecase.ProposalEnglishItemUsecase, createUse usecase.CreateEnglishItemUsecase, getUse usecase.GetEnglishItemUsecase, updateUse usecase.UpdateEnglishItemUsecase, deleteUse usecase.DeleteEnglishItemUsecase, getRequiredExpUse usecase.GetRequiredExpUsecase) EnglishItemController {
	return &EnglishItemGinController{
		proposalUse:       proposalUse,
		createUse:         createUse,
		getUse:            getUse,
		updateUse:         updateUse,
		deleteUse:         deleteUse,
		getRequiredExpUse: getRequiredExpUse,
	}
}

func (ec *EnglishItemGinController) Proposal(c *gin.Context) {
	req := &dto.ProposalEnglishItemRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	resp, err := ec.proposalUse.Proposal(req)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) ProposalTranslation(c *gin.Context) {
	req := &dto.ProposalEnglishItemRequest{}

	if err := c.ShouldBindQuery(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	resp, err := ec.proposalUse.ProposalTranslation(req)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.String(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) ProposalExplanation(c *gin.Context) {
	req := &dto.ProposalEnglishItemRequest{}

	if err := c.ShouldBindQuery(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	resp, err := ec.proposalUse.ProposalExplanation(req)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.String(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) ProposalExample(c *gin.Context) {
	req := &dto.ProposalEnglishItemRequest{}

	if err := c.ShouldBindQuery(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	resp, err := ec.proposalUse.ProposalExample(req)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) Create(c *gin.Context) {
	req := &dto.CreateEnglishItemRequest{}
	if err := c.BindJSON(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
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

func (ec *EnglishItemGinController) Update(c *gin.Context) {
	req := &dto.UpdateEnglishItemRequest{}
	if err := c.BindJSON(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	userId, err := userId(c)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	req.UserId = userId
	if err := ec.updateUse.Update(req); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.Status(http.StatusCreated)
}

func (ec *EnglishItemGinController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	if err := ec.deleteUse.Delete(id); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func (ec *EnglishItemGinController) GetByUserId(c *gin.Context) {
	userId, err := userId(c)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	resp, err := ec.getUse.GetEnglishItemInfoByUserId(userId)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
	}

	if resp == nil {
		resp = &dto.GetEnglishItemsResponse{
			EnglishItems: []*dto.EnglishItem{},
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) GetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	resp, err := ec.getUse.GetById(id)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ec *EnglishItemGinController) GetRequiredExp(c *gin.Context) {
	output := ec.getRequiredExpUse.GetRequiredExp()

	c.JSON(http.StatusOK, output)
}
