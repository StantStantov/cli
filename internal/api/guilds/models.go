package guilds

const (
	JoinGuildPath        = "/guilds/%d/join"            // POST - Вступление в гильдию
	CreateGuildPath      = "/guilds"                    // POST - Создание гильдии
	DeleteGuildPath      = "/guilds/%d"                 // DELETE - Удаление гильдии
	GetGuildMembersPath  = "/guilds/%d/members"         // GET - Список участников
	UpdateMemberRolePath = "/guilds/%d/members/%d/role" // PUT - Изменение роли
)

// BaseResponse - базовый формат ответа
type BaseResponse struct {
	Error     string      `json:"error"`
	ErrorCode int         `json:"error_code"`
	Status    bool        `json:"status"`
	Value     interface{} `json:"value"`
}

// GuildResponse - данные гильдии
type GuildResponse struct {
	ID          int    `json:"id"`
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     int    `json:"owner_id"`
}

// GuildMember - участник гильдии
type GuildMember struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

// CreateGuildRequest - запрос на создание гильдии
type CreateGuildRequest struct {
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateRoleRequest - запрос на изменение роли
type UpdateRoleRequest struct {
	RoleID int `json:"role_id"`
}

// ErrorResponse - ответ об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}
