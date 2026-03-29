package pkg

import (
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	pt_br_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

func BadRequestResponseParser() gin.HandlerFunc {

	universalTranslator := registerTranslations()

	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() != http.StatusBadRequest {
			return
		}

		errs := c.Errors.Last()
		if errs == nil {
			return
		}

		ve, ok := errs.Err.(validator.ValidationErrors)
		if !ok {
			return
		}

		language := c.GetHeader("Accept-Language")
		language = strings.ReplaceAll(language, "-", "_")
		if language == "" {
			language = "en"
		}

		translator, _ := universalTranslator.GetTranslator(language)

		translatedErrors := ve.Translate(translator)
		c.JSON(http.StatusBadRequest, translatedErrors)
	}
}

func registerTranslations() *ut.UniversalTranslator {
	ptBrT := pt_BR.New()
	enT := en.New()
	
	universalTranslator := ut.New(enT, enT, ptBrT) 
	engine, _ := binding.Validator.Engine().(*validator.Validate)

	enTrans, _ := universalTranslator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(engine, enTrans)

	ptBrTrans, _ := universalTranslator.GetTranslator("pt_BR")
	pt_br_translations.RegisterDefaultTranslations(engine, ptBrTrans)

	return universalTranslator
}
