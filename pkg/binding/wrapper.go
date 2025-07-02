package binding

import (
	"encoding/json"
	"fmt"
	logger "go-source/pkg/log"
	"go-source/pkg/resp"
	validator "go-source/pkg/validate"
	"net/http"

	"github.com/labstack/echo/v4"
)

const RequestObjectContextKey = "service_requestObject"

// Simple error response structure
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Lang    string `json:"lang"`
}

type Error struct {
	CustomCode int `json:"code,omitempty" example:"400"`

	*echo.HTTPError
}

func (e *Error) Error() string {
	return fmt.Sprintf("CustomCode=%d, %s", e.CustomCode, e.HTTPError.Error())
}

func ErrorHandler(next echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		// same check as default handler
		if c.Response().Committed {
			return
		}

		if e, ok := err.(*Error); ok {
			c.JSON(e.Code, resp.BuildErrorResp(resp.ErrDataInvalid, e.Error(), resp.LangEN))
			return
		}

		if next != nil {
			next(err, c)
		}
	}
}

func RestLogFieldExtractor(c echo.Context) []interface{} {
	if req := c.Get(RequestObjectContextKey); req != nil {
		reqObjectString := ""
		if b, err := json.Marshal(req); err != nil {
			reqObjectString = fmt.Sprintf("fail to parse reqObject as string: %+v", err)
		} else {
			reqObjectString = string(b)
		}

		return []interface{}{"requestObject", reqObjectString}
	}

	return nil
}

func Wrapper[TReq any, TRes any](wrapped func(echo.Context, *TReq) (*TRes, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		logging := logger.GetLogger()

		var req TReq
		if err := c.Bind(&req); err != nil {
			logging.Error().Err(err).Str("request_uri", c.Request().RequestURI).Msg("fail to bind request")
			return c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		}

		// Use custom validator instead of Echo's built-in validation
		v := validator.NewValidate()
		if err := v.ValidateStruct(&req); err != nil {
			logging.Error().Err(err).Str("request_uri", c.Request().RequestURI).Str("request_object", fmt.Sprintf("%+v", req)).Msg("fail to validate request")
			return c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		}

		c.Set(RequestObjectContextKey, req)

		// Let the handler manage its own response format
		_, err := wrapped(c, &req)
		return err
	}
}
