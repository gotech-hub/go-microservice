package models

type GetByProfileIdRequest struct {
	ProfileID string `json:"profileId" validate:"required"`
}
