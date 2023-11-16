package controller

import (
	"english/cmd/presentation/errhandle"
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"english/myerror"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OutputController interface {
	GetQuestion(c *gin.Context)
	AnswerQuestions(c *gin.Context)
	CreateOutput(c *gin.Context)
	GetOutputTimes(c *gin.Context)
	GetOutputs(c *gin.Context)
	DeleteOutput(c *gin.Context)
}

type OutputGinController struct {
	getQuestionUse     usecase.GetQuestionUsecase
	answerQuestionsUse usecase.AnswerQuestionsUsecase
	createOutputUse    usecase.CreateOutputUsecase
	getOutputTimesUse  usecase.GetOutputTimesUsecase
	getOutputsUse      usecase.GetOutputsUsecase
	deleteOutputUse    usecase.DeleteOutputUsecase
}

func NewOutputController(getQuestionUse usecase.GetQuestionUsecase, answerQuestionsUse usecase.AnswerQuestionsUsecase, createOutputUse usecase.CreateOutputUsecase, getOutputTimesUse usecase.GetOutputTimesUsecase, getOutputsUse usecase.GetOutputsUsecase, deleteOutputUse usecase.DeleteOutputUsecase) OutputController {
	return &OutputGinController{
		getQuestionUse:     getQuestionUse,
		answerQuestionsUse: answerQuestionsUse,
		createOutputUse:    createOutputUse,
		getOutputTimesUse:  getOutputTimesUse,
		getOutputsUse:      getOutputsUse,
		deleteOutputUse:    deleteOutputUse,
	}
}

func (oc OutputGinController) GetQuestion(c *gin.Context) {
	input := &dto.GetQuestionInput{}
	if err := c.BindQuery(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%w:%v", myerror.ErrBindingFailure, err), c)
		return
	}

	if err := Validate(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%w:%v", myerror.ErrValidation, err), c)
		return
	}

	output, err := oc.getQuestionUse.GetQuestion(input)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (oc *OutputGinController) AnswerQuestions(c *gin.Context) {
	input := &dto.AnswerQuestionsInput{}
	if err := c.BindJSON(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%v: %w", myerror.ErrBindingFailure, err), c)
		return
	}

	if err := Validate(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%v: %w", myerror.ErrValidation, err), c)
		return
	}

	output, err := oc.answerQuestionsUse.AnswerQuestions(input)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (oc *OutputGinController) CreateOutput(c *gin.Context) {
	input := &dto.CreateOutputInput{}

	if err := c.BindJSON(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%v: %w", myerror.ErrBindingFailure, err), c)
		return
	}

	if err := Validate(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%v: %w", myerror.ErrValidation, err), c)
		return
	}

	if err := oc.createOutputUse.Create(input); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.Status(http.StatusCreated)
}

func (oc *OutputGinController) GetOutputTimes(c *gin.Context) {
	englishItemId := c.Param("englishItemId")

	if englishItemId == "" {
		errhandle.HandleErrorJSON(myerror.ErrValidation, c)
		return
	}

	output, err := oc.getOutputTimesUse.GetOutputTimes(englishItemId)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (oc *OutputGinController) GetOutputs(c *gin.Context) {
	input := &dto.GetOutputsInput{}

	if err := c.BindQuery(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%w:%v", myerror.ErrBindingFailure, err), c)
		return
	}

	if err := Validate(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%w:%v", myerror.ErrValidation, err), c)
		return
	}

	output, err := oc.getOutputsUse.GetOutputs(input)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (oc *OutputGinController) DeleteOutput(c *gin.Context) {
	input := &dto.DeleteOutputInput{}

	if err := c.BindQuery(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%w:%v", myerror.ErrBindingFailure, err), c)
		return
	}

	if err := Validate(input); err != nil {
		errhandle.HandleErrorJSON(fmt.Errorf("%w:%v", myerror.ErrValidation, err), c)
		return
	}

	if err := oc.deleteOutputUse.Delete(input); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.Status(http.StatusOK)
}
