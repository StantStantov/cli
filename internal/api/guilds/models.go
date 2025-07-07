package guilds

// константы путей API
const (
	PathGetMemberByUserID = "/api/v1/guild/member/member/%d"
	PathGetGuildByTag     = "/api/v1/guild/%s"
	PathSendJoinRequest   = "/api/v1/guild/request/%s"
	PathGetJoinRequests   = "/api/v1/guild/request/%s"
	PathApplyJoinRequest  = "/api/v1/guild/request/%s/%d/apply"
	PathCancelJoinRequest = "/api/v1/guild/request/%s/%d/cancel"
	PathCreateGuild       = "/api/v1/guild/"
	PathDeleteGuild       = "/api/v1/guild/%s"
	PathGetGuildMembers   = "/api/v1/guild/member/%s"
	PathDeleteMember      = "/api/v1/guild/member/%s/%d"
	PathEditMember        = "/api/v1/guild/member/%s/%d"
	PathExitGuild         = "/api/v1/guild/member/%s/exit/%d"
)

// Role - роль участника гильдии
type Role struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	RolePromote []int  `json:"role_promote"` // список ID ролей, которыми может управлять пользователь
	// permissions опущены
}

// MemberResponse - информация об участнике гильдии
type MemberResponse struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	GuildID  int    `json:"guild_id"`
	GuildTag string `json:"guild_tag"`
	Role     Role   `json:"role"`
}

// GuildResponse - информация о гильдии
type GuildResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	OwnerID     int    `json:"owner_id"`
	IsActive    bool   `json:"is_active"` // активна ли
	IsFull      bool   `json:"is_full"`   // заполнена ли
}

// GuildPagination - список гильдий с пагинацией
type GuildPagination struct {
	Items      []GuildResponse `json:"items"`
	TotalItems int             `json:"total_items"`
	TotalPages int             `json:"total_pages"`
}

// MemberPagination - список участников с пагинацией
type MemberPagination struct {
	Items      []MemberResponse `json:"items"`
	TotalItems int              `json:"total_items"`
	TotalPages int              `json:"total_pages"`
}

// BaseResponse - базовый формат ответа
type BaseResponse struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Status    bool   `json:"status"`
}

// ResponseGuild - ответ с данными гильдии
type ResponseGuild struct {
	BaseResponse
	Value *GuildResponse `json:"value"`
}

// ResponseGuildPagination - ответ со списком гильдий
type ResponseGuildPagination struct {
	BaseResponse
	Value *GuildPagination `json:"value"`
}

// ResponseMember - ответ с данными участника
type ResponseMember struct {
	BaseResponse
	Value *MemberResponse `json:"value"`
}

// ResponseMemberPagination - ответ со списком участников
type ResponseMemberPagination struct {
	BaseResponse
	Value *MemberPagination `json:"value"`
}

// RequestResponse - заявка на вступление
type RequestResponse struct {
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// RequestPagination - список заявок с пагинацией
type RequestPagination struct {
	Items      []RequestResponse `json:"items"`
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
}

// ResponseRequestPagination - ответ со списком заявок
type ResponseRequestPagination struct {
	BaseResponse
	Value *RequestPagination `json:"value"`
}

// CreateGuildRequest - запрос на создание гильдии
type CreateGuildRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

// EditMemberRequest - запрос на изменение участника
type EditMemberRequest struct {
	RoleID   int    `json:"role_id"`   // ID новой роли
	UserName string `json:"user_name"` // Новое имя (опционально)
}
