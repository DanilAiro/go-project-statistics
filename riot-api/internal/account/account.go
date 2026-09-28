package account

import (
	"strconv"
	"strings"
)

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

func (a Account) String() string {
	var b strings.Builder
	b.WriteString("Tag: ")
	b.WriteString(a.Tag)
	b.WriteByte('\n')
	b.WriteString("Name: ")
	b.WriteString(a.Name)
	b.WriteByte('\n')
	b.WriteString("Card: ")
	b.WriteString(a.Card)
	b.WriteByte('\n')
	b.WriteString("PUUID: ")
	b.WriteString(a.PUUID)
	b.WriteByte('\n')
	b.WriteString("Title: ")
	b.WriteString(a.Title)
	b.WriteByte('\n')
	b.WriteString("Region: ")
	b.WriteString(a.Region)
	b.WriteByte('\n')
	b.WriteString("UpdatedAt: ")
	b.WriteString(a.UpdatedAt)
	b.WriteByte('\n')
	b.WriteString("Platforms: ")
	b.WriteString("[" + strings.Join(a.Platforms, ", ") + "]")
	b.WriteByte('\n')
	b.WriteString("AccountLevel: ")
	b.WriteString(strconv.FormatInt(int64(a.AccountLevel), 10))

	return b.String()
}
