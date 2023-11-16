package main

import (
	"english/cmd/infrastructure/client"
	"english/cmd/infrastructure/dbconn"
	"english/cmd/infrastructure/gateway"
	"english/cmd/presentation/controller"
	"english/cmd/presentation/router"
	"english/cmd/usecase"
	"english/cmd/usecase/validator"
	"english/config"
	"fmt"
)

func init() {
	config.SetLogger(config.LogFileName())
}

func main() {
	db := dbconn.NewDB()

	gcsClient := client.NewGCSClient()
	fileStorageRepo := gateway.NewFileStorageGCSRepository(gcsClient)

	userRepo := gateway.NewUserMySQLRepository(db)
	signupUse := usecase.NewSignupUsecase(userRepo)
	loginUse := usecase.NewLoginUsecase(userRepo)
	ssoAuthUse := usecase.NewSSOAuthUsecase(userRepo)
	getUserUse := usecase.NewGetUserUsecase(userRepo)
	updateUserUse := usecase.NewUpdateUserProfileUsecase(userRepo)
	updateProfileUse := usecase.NewUpdateProfileImgUsecase(userRepo, fileStorageRepo)
	uc := controller.NewUserGinController(signupUse, loginUse, ssoAuthUse, getUserUse, updateUserUse, updateProfileUse)

	chatGPTAPI := client.NewOpenAIAPI()
	proposalUse := usecase.NewProposalEnglishItemUsecase(chatGPTAPI)
	englishItemRepo := gateway.NewEnglishItemMySQLReporitory(db)
	englishItemValidator := validator.NewEnglishItemValidator()
	createEnglishItemUse := usecase.NewCreateEnglishItemUsecase(englishItemRepo, fileStorageRepo, englishItemValidator)
	findAllEnglishItemUse := usecase.NewGetEnglishItemUsecase(englishItemRepo)
	updateEnglishItemUse := usecase.NewUpdateEnglishItemUsecase(englishItemRepo, fileStorageRepo, englishItemValidator)
	deleteEnglishItemUse := usecase.NewDeleteEnglishItemUsecase(englishItemRepo, fileStorageRepo)
	getRequiredExpUse := usecase.NewGetRequiredExpUsecase()
	ec := controller.NewEnglishItemController(proposalUse, createEnglishItemUse, findAllEnglishItemUse, updateEnglishItemUse, deleteEnglishItemUse, getRequiredExpUse)

	outputRepo := gateway.NewOutputRepository(db)
	getQuestionUse := usecase.NewGetQuestionUsecase(chatGPTAPI)
	answerQuestionsUse := usecase.NewAnswerQuestionsUsecase(englishItemRepo, chatGPTAPI)
	createOutputUse := usecase.NewCreateOutputUsecase(outputRepo)
	getOutputTimesUse := usecase.NewGetOutputTimesUsecase(outputRepo)
	getOutputsUse := usecase.NewGetOutputsUsecase(outputRepo)
	deleteOutputUse := usecase.NewDeleteOutputUsecase(outputRepo)
	oc := controller.NewOutputController(getQuestionUse, answerQuestionsUse, createOutputUse, getOutputTimesUse, getOutputsUse, deleteOutputUse)

	router := router.NewGinRouter(uc, ec, oc)

	router.Run(fmt.Sprintf(":%v", config.Port()))
}
