package workspace

import "time"

type WorkspaceStatus string
type WorkspaceType string

const (
	Active  WorkspaceStatus = "active"
	Deleted WorkspaceStatus = "deleted"
)

const (
	Advertiser WorkspaceType = "advertiser"
	Publisher  WorkspaceType = "publisher"
)

type Workspace struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Type        WorkspaceType   `json:"type"`
	Country     string          `json:"country"`
	Timezone    string          `json:"timezone"`
	Currency    string          `json:"currency"`
	Industry    string          `json:"industry"`
	Website     string          `json:"website"`
	Description string          `json:"description"`
	Status      WorkspaceStatus `json:"status"`
	CreatedBy   string          `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateWorkspaceDTO struct {
	Name        string        `json:"name" validate:"required,max=150"`
	Type        WorkspaceType `json:"type" validate:"required,oneof=advertiser publisher"`
	Country     string        `json:"country" validate:"required,len=2"`
	Timezone    string        `json:"timezone" validate:"required,max=100"`
	Currency    string        `json:"currency" validate:"required,len=3"`
	Industry    *string       `json:"industry" validate:"omitempty,max=100"`
	Website     *string       `json:"website" validate:"omitempty,url"`
	Description *string       `json:"description"`
}

type UpdateWorkspaceDTO struct {
	Name        *string `json:"name"`
	Country     *string `json:"country"`
	Timezone    *string `json:"timezone"`
	Currency    *string `json:"currency"`
	Industry    *string `json:"industry"`
	Website     *string `json:"website"`
	Description *string `json:"description"`
}
