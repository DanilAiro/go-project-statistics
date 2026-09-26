package match

type Match struct {
	Kills     []kill     `json:"kills"`
	Teams     []team     `json:"teams"`
	Rounds    []round    `json:"rounds"`
	Players   []player   `json:"players"`
	Coaches   []coache   `json:"coaches"`
	Metadata  metadata   `json:"metadata"`
	Observers []observer `json:"observers"`
}

type coache struct {
	PUUID  string `json:"puuid"`
	TeamID string `json:"team_id"`
}

type kill struct {
	Round             int32            `json:"round"`
	Killer            killer           `json:"killer"`
	Victim            victim           `json:"victim"`
	Weapon            weapon           `json:"weapon"`
	Location          location         `json:"location"`
	Assistants        []assistant      `json:"assistants"`
	PlayerLocations   []playerLocation `json:"player_locations"`
	TimeInRoundInMS   int64            `json:"time_in_round_in_ms"`
	TimeInMatchInMS   int64            `json:"time_in_match_in_ms"`
	SecondaryFireMode bool             `json:"secondary_fire_mode"`
}

type metadata struct {
	Map             gameMap           `json:"map"`
	Queue           queue             `json:"queue"`
	Region          string            `json:"region"`
	MatchID         string            `json:"match_id"`
	Cluster         string            `json:"cluster"`
	Premier         any               `json:"premier"`
	Platform        string            `json:"platform"`
	StartedAt       string            `json:"started_at"`
	GameVersion     string            `json:"game_version"`
	IsCompleted     bool              `json:"is_completed"`
	GameLengthInMS  int64             `json:"game_length_in_ms"`
	PartyRRPenaltys []partyRRPenaltys `json:"party_rr_penaltys"`
}

type observer struct {
	Tag                 string `json:"tag"`
	Name                string `json:"name"`
	CardID              string `json:"card_id"`
	PUUID               string `json:"puuid"`
	TitleID             string `json:"title_id"`
	PartyID             string `json:"party_id"`
	AccountLevel        int32  `json:"account_level"`
	SessionPlaytimeInMS int32  `json:"session_playtime_in_ms"`
}

type player struct {
	Tag                 string        `json:"tag"`
	Name                string        `json:"name"`
	Tier                tier          `json:"tier"`
	Agent               agent         `json:"agent"`
	PUUID               string        `json:"puuid"`
	Stats               stats         `json:"stats"`
	TeamID              string        `json:"team_id"`
	Economy             economy       `json:"economy"`
	PartyID             string        `json:"party_id"`
	Behavior            behavior      `json:"behavior"`
	Platform            string        `json:"platform"`
	AbilityCasts        abilityCasts  `json:"ability_casts"`
	AccountLevel        int32         `json:"account_level"`
	Customization       customization `json:"customization"`
	SessionPlaytimeInMS int32         `json:"session_playtime_in_ms"`
}

type round struct {
	ID          int32   `json:"id"`
	Stats       []stats `json:"stats"`
	Plant       plant   `json:"plant"`
	Result      string  `json:"result"`
	Defuse      defuse  `json:"defuse"`
	Ceremony    string  `json:"ceremony"`
	WinningTeam string  `json:"winning_team"`
}

type team struct {
	Won           bool          `json:"won"`
	Rounds        []roundShort  `json:"rounds"`
	TeamID        string        `json:"team_id"`
	PremierRoster premierRoster `json:"premier_roster"`
}

type assistant struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	Team  string `json:"team"`
	PUUID string `json:"puuid"`
}

type killer struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	Team  string `json:"team"`
	PUUID string `json:"puuid"`
}

type location struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type playerLocation struct {
	Player      playerShort `json:"player"`
	Location    location    `json:"location"`
	ViewRadians float64     `json:"view_radians"`
}

type playerShort struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	Team  string `json:"team"`
	PUUID string `json:"puuid"`
}

type victim struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	Team  string `json:"team"`
	PUUID string `json:"puuid"`
}

type weapon struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type gameMap struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type partyRRPenaltys struct {
	PartyID string  `json:"party_id"`
	Penalty float64 `json:"penalty"`
}

type queue struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ModeType string `json:"mode_type"`
}

type season struct {
	ID    string `json:"id"`
	Short string `json:"short"`
}

type abilityCasts struct {
	Ability1 int32 `json:"abitity1"`
	Ability2 int32 `json:"abitity2"`
	Grenade  int32 `json:"grenade"`
	Ultimate int32 `json:"ultimate"`
}

type agent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type behavior struct {
	AfkRounds     float64      `json:"afk_rounds"`
	FriendlyFire  friendlyFire `json:"friendly_fire"`
	RoundsInSpawn float64      `json:"rounds_in_spawn"`
}

type friendlyFire struct {
	Incoming float64 `json:"incoming"`
	Outgoing float64 `json:"outgoing"`
}

type customization struct {
	Card                 string `json:"card"`
	Title                string `json:"title"`
	PreferredLevelBorder string `json:"preferred_level_border"`
}

type economy struct {
	Spent        spent        `json:"spent"`
	LoadoutValue loadoutValue `json:"loadout_value"`
}

type loadoutValue struct {
	Overall int32   `json:"overall"`
	Average float64 `json:"average"`
}

type spent struct {
	Overall int32   `json:"overall"`
	Average float64 `json:"average"`
}

type stats struct {
	Kills     int32  `json:"kills"`
	Score     int32  `json:"score"`
	Damage    damage `json:"damage"`
	Assists   int32  `json:"assists"`
	Legshots  int32  `json:"legshots"`
	Bodyshots int32  `json:"bodyshots"`
	Headshots int32  `json:"headshots"`
}

type damage struct {
	Dealt    int32 `json:"dealt"`
	Received int32 `json:"received"`
}

type tier struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type defuse struct {
	Player          playerShort      `json:"player"`
	Location        location         `json:"location"`
	RoundTimeInMS   int32            `json:"round_time_in_ms"`
	PlayerLocations []playerLocation `json:"player_locations"`
}

type plant struct {
	Player          playerShort      `json:"player"`
	Location        location         `json:"location"`
	RoundTimeInMS   int32            `json:"round_time_in_ms"`
	PlayerLocations []playerLocation `json:"player_locations"`
}

type roundShort struct {
	Won  int32 `json:"won"`
	Lost int32 `json:"lost"`
}

type premierRoster struct {
	ID            string               `json:"id"`
	Tag           string               `json:"tag"`
	Name          string               `json:"name"`
	Members       []string             `json:"members"`
	Customization premierCustomization `json:"customization"`
}

type premierCustomization struct {
	Icon           string `json:"icon"`
	Image          string `json:"image"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	TetritaryColor string `json:"tetritary_color"`
}
