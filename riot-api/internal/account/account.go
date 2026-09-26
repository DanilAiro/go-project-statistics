package account

type Account struct {
	Tag          string   `json:"tag"`
	Name         string   `json:"name"`
	Card         string   `json:"card"`
	PUUID        string   `json:"puuid"`
	Title        string   `json:"title"`
	Region       string   `json:"region"`
	UpdatedAt    string   `json:"updated_at"`
	Platforms    []string `json:"platforms"`
	AccountLevel int32    `json:"account_level"`
}
