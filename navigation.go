package betfair

type NavigationChild struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Children []NavigationChild
}

type Navigation struct {
	Children []NavigationChild `json:"children"`
}
