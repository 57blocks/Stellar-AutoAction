package model

type UserInfo struct {
	Id             string `json:"id"`
	Account        string `json:"account"`
	Password       string `json:"password"`
	Description    string `json:"description"`
	OrganizationId int32  `json:"organization_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type OrganizationInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TokenInfo struct {
	Access         string `json:"access"`
	Refresh        string `json:"refresh"`
	UserId         int32  `json:"user_id"`
	AccessExpires  string `json:"access_expires"`
	RefreshExpires string `json:"refresh_expires"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
