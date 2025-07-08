package guild

import (
	"lesta-start-battleship/cli/internal/api/guilds"
	"sync"
)

var (
	mu sync.RWMutex

	// Список участников гильдии: username -> MemberResponse
	Members = make(map[string]guilds.MemberResponse)

	// Информация о текущем пользователе в гильдии
	Self guilds.MemberResponse

	// Список гильдий: guild_tag -> GuildResponse
	GuildList = make(map[string]guilds.GuildResponse)
)

func SetMember(username string, member guilds.MemberResponse) {
	mu.Lock()
	defer mu.Unlock()
	Members[username] = member
}

func GetMember(username string) (guilds.MemberResponse, bool) {
	mu.RLock()
	defer mu.RUnlock()
	val, ok := Members[username]
	return val, ok
}

func SetSelf(member guilds.MemberResponse) {
	mu.Lock()
	defer mu.Unlock()
	Self = member
}

func SetGuild(guildTag string, guild guilds.GuildResponse) {
	mu.Lock()
	defer mu.Unlock()
	GuildList[guildTag] = guild
}

func GetGuild(guildTag string) (guilds.GuildResponse, bool) {
	mu.RLock()
	defer mu.RUnlock()
	val, ok := GuildList[guildTag]
	return val, ok
}
