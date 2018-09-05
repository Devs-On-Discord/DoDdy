package roles

type Role int

var (
	roleNames = [...]string{
		"User",
		"Node Mod",
		"Node Admin",
		"Hub Mod",
		"Hub Admin",
		"Bot Developer",
	}
	CommandRoleNames = map[string]Role{
		"User":         User,
		"NodeMod":      NodeMod,
		"NodeAdmin":    NodeAdmin,
		"HubMod":       HubMod,
		"HubAdmin":     HubAdmin,
		"BotDeveloper": BotDeveloper,
	}
	RoleInt = map[int]Role{
		0: User,
		1: NodeMod,
		2: NodeAdmin,
		3: HubMod,
		4: HubAdmin,
		5: BotDeveloper,
	}
)

const (
	User         Role = 0
	NodeMod      Role = 1
	NodeAdmin    Role = 2
	HubMod       Role = 3
	HubAdmin     Role = 4
	BotDeveloper Role = 5
)

func (role Role) String() string {
	if role < User || role > BotDeveloper {
		return "Unknown"
	}
	return roleNames[role]
}
