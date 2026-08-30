package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// FieldError describes a single invalid field in a request body.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// BindError writes a 400 response identifying which field(s) in the request
// body are missing or invalid, in place of a bare "Invalid Request Body".
// Pass it the error returned by c.ShouldBindJSON.
func BindError(c *gin.Context, err error) {
	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		fields := make([]FieldError, 0, len(ve))
		for _, fe := range ve {
			fields = append(fields, FieldError{
				Field:   fe.Field(),
				Message: validationMessage(fe),
			})
		}
		respondInvalidBody(c, fields)
		return
	}

	if ute, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		respondInvalidBody(c, []FieldError{{
			Field:   ute.Field,
			Message: fmt.Sprintf("must be a %s", ute.Type.String()),
		}})
		return
	}

	if errors.Is(err, io.EOF) {
		respondInvalidBody(c, []FieldError{{
			Field:   "body",
			Message: "request body is required",
		}})
		return
	}

	respondInvalidBody(c, []FieldError{{
		Field:   "body",
		Message: err.Error(),
	}})
}

func respondInvalidBody(c *gin.Context, fields []FieldError) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": MsgInvalidBody,
		"errors":  fields,
	})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", fe.Param())
	default:
		return fmt.Sprintf("failed validation on %q", fe.Tag())
	}
}
