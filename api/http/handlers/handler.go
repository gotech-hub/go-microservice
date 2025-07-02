package handlers

import (
	"go-source/api/http/models"
	"go-source/pkg/resp"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewHandler() (handler *Handler) {
	return &Handler{}
}

// GetByProfileId godoc
// @Summary Get by profile ID
// @Description Get all for a specific profile ID
// @Tags tier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId header string true "Profile ID"
// @Param clientId query string true "Client ID"
// @Param X-Client-Region header string false "Client Region (e.g., VN, TH, ID, MY, SG, PH)"
// @Success 200 {object} resp.Resp
// @Failure 400 {object} resp.Resp
// @Failure 404 {object} resp.Resp
// @Router /v1/profile [get]
func (h *Handler) GetByProfileId(c echo.Context, req *models.GetByProfileIdRequest) (*resp.Resp, error) {
	ctx := c.Request().Context()
	profileID, ok := ctx.Value("").(string)
	if !ok {
		return nil, c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, "", resp.LangEN))
	}
	if profileID == "" {
		return nil, c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, "", resp.LangEN))
	}

	// tiers, err := h.tierService.GetTiersByProfileId(ctx, profileID, clientID)
	// if err != nil {
	// 	return c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
	// }

	// return c.JSON(http.StatusOK, resp.BuildSuccessResp(resp.LangEN, tiers))
	return resp.BuildSuccessResp(resp.LangEN, "GetByProfileId handler not implemented yet"), nil
}
