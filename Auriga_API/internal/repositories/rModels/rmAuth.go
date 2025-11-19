package rModels

/*
type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
} */

// AuthentikTokenResponse representa la respuesta de token de Authentik
type AuthentikTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// AuthentikUserInfo representa la información del usuario de Authentik
type AuthentikUserInfo struct {
	Sub               string                 `json:"sub"`
	Email             string                 `json:"email"`
	Name              string                 `json:"name"`
	PreferredUsername string                 `json:"preferred_username"`
	Groups            []string               `json:"groups"`
	Organization      map[string]interface{} `json:"organization"`
}

// UserOrganizationInfo representa la información de organización del usuario
type UserOrganizationInfo struct {
	OrganizationName string              `json:"organization_name"`
	Factories        []UserFactoryAccess `json:"factories"`
}

type UserFactoryAccess struct {
	Factory     string   `json:"factory"`
	Department  string   `json:"department"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
