package resp

const (
	ErrDataInvalid = 1002
	LangEN         = "en"
	LangVI         = "vi"
)

type Resp struct {
	ErrorCode   interface{} `json:"errorCode"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	Lang        string      `json:"lang"`
}

func BuildErrorResp(errorCode interface{}, description string, lang string) *Resp {
	return &Resp{
		ErrorCode:   errorCode,
		Message:     "Error",
		Description: description,
		Lang:        lang,
	}
}

func BuildSuccessResp(lang string, data interface{}) *Resp {
	return &Resp{
		ErrorCode:   0,
		Message:     "Success",
		Description: "",
		Lang:        lang,
	}
}
