package http

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	pt_br_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/jeffersonbrasilino/ddgo"
)

var GlobalTranslator *ut.UniversalTranslator

func init() {
	GlobalTranslator = registerTranslations()
}

func registerTranslations() *ut.UniversalTranslator {
	ptBrT := pt_BR.New()
	enT := en.New()

	universalTranslator := ut.New(enT, enT, ptBrT)

	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		enTrans, _ := universalTranslator.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(engine, enTrans)

		ptBrTrans, _ := universalTranslator.GetTranslator("pt_BR")
		pt_br_translations.RegisterDefaultTranslations(engine, ptBrTrans)
	}

	return universalTranslator
}

func Error(c *gin.Context, err error) {
	switch err.(type) {
	case *ddgo.ValidationError:
		ErrorWithCode(c, 400, err)
	case *ddgo.NotFoundError:
		ErrorWithCode(c, 404, err)
	case *ddgo.AlreadyExistsError:
		ErrorWithCode(c, 409, err)
	case *ddgo.DependencyError:
		ErrorWithCode(c, 502, err)
	case *ddgo.InvalidDataError:
		ErrorWithCode(c, 422, err)
	default:
		ErrorWithCode(c, 500, err)
	}
}

func ErrorWithCode(c *gin.Context, code int, err error) {

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		message := err.Error()
		var raw json.RawMessage
		if errJson := json.Unmarshal([]byte(message), &raw); errJson == nil {
			c.JSON(code, gin.H{
				"errors": raw,
			})
			return
		}

		c.JSON(code, gin.H{
			"errors": message,
		})

		return
	}

	language := strings.Split(c.GetHeader("Accept-Language"), ",")[0]
	language = strings.ReplaceAll(language, "-", "_")
	if language == "" {
		language = "en"
	}

	translator, _ := GlobalTranslator.GetTranslator(language)
	translatedErrors := ve.Translate(translator)

	c.JSON(code, gin.H{
		"errors": translatedErrors,
	})
}

func Success(c *gin.Context, code int, data any) {
	dataResponse, _ := json.Marshal(data)
	c.Data(code, gin.MIMEJSON, dataResponse)
}
