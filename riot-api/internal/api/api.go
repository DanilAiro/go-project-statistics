package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// GetAccountByName возвращает данные пользователя VALORANT по имени и тегу.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - name — имя игрока (Riot ID)
//   - tag — тег игрока (Riot ID tag line)
//
// query содержит GET query-параметры, которые будут добавлены к URL запроса (опционально):
//   - force bool — обновить данные или забрать из кеша
func GetAccountByName(apiKey string, params map[string]string, query map[string]any) {
	body := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v2/account", params, query)

	fmt.Println(string(body))
}

// GetAccount возвращает данные пользователя VALORANT по идентификатору пользователя.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - puuid — идентификатор игрока
//
// query содержит GET query-параметры, которые будут добавлены к URL запроса (опционально):
//   - force bool — обновить данные или забрать из кеша
func GetAccount(apiKey string, params map[string]string, query map[string]any) {
	body := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v2/by-puuid/account", params, query)

	fmt.Println(string(body))
}

// GetPlayerMMR возвращает рейтинг пользователя VALORANT по идентификатору пользователя.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - affinity — регион игрока (e.g., na, eu, ap, kr)
//   - platform — платформа игрока (pc, console)
//   - puuid — идентификатор игрока
func GetPlayerMMR(apiKey string, params map[string]string) {
	body := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v3/by-puuid/mmr", params, nil)

	fmt.Println(string(body))
}

// GetPlayerMatches возвращает матчи пользователя VALORANT по идентификатору пользователя.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - affinity — регион игрока (e.g., na, eu, ap, kr)
//   - platform — платформа игрока (pc, console)
//   - puuid — идентификатор игрока
//
// query содержит GET query-параметры, которые будут добавлены к URL запроса (опционально):
//   - mode string — игровой режим
//   - map string — название карты
//   - size int32 — количество результатов
//   - start int32 — начальный индекс для пагинации результатов
func GetPlayerMatches(apiKey string, params map[string]string, query map[string]any) {
	body := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v3/by-puuid/matches", params, query)

	fmt.Println(string(body))
}

// GetMatchInfo возвращает данные матча VALORANT по идентификатору матча.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - affinity — регион игрока (e.g., na, eu, ap, kr)
//   - match_id — идентификатор матча
func GetMatchInfo(apiKey string, params map[string]string) {
	body := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v4/match", params, nil)

	fmt.Println(string(body))
}

// Вспомогательные функции
func getEndpoint(path string, params map[string]string) string {
	endpoint := path

	for _, v := range params {
		endpoint = fmt.Sprintf(
			"%s/%s",
			endpoint,
			url.PathEscape(v),
		)
	}

	return endpoint
}

func setQuery(urlValues url.Values, query map[string]any) {
	for k, v := range query {
		if v == nil {
			continue // или urlValues.Set(k, ""), если нужно передать пустой ключ
		}

		var strVal string
		switch val := v.(type) {
		case string:
			strVal = val
		case int:
			strVal = strconv.Itoa(val)
		case int64:
			strVal = strconv.FormatInt(val, 10)
		case bool:
			strVal = strconv.FormatBool(val)
		case float64:
			// 'f' убирает экспоненциальную нотацию (например, 1e+06 -> 1000000)
			strVal = strconv.FormatFloat(val, 'f', -1, 64)
		default:
			// Для сложных типов (массивы, структуры, слайсы) используем форматирование fmt
			strVal = fmt.Sprintf("%v", val)
		}

		urlValues.Set(k, strVal)
	}
}

func baseRequest(apiKey, url string, params map[string]string, query map[string]any) []byte {
	endpoint := getEndpoint(url, params)

	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)

	if query != nil {
		urlValues := req.URL.Query()
		setQuery(urlValues, query)
		req.URL.RawQuery = urlValues.Encode()
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body
}
