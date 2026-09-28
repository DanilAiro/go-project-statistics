package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"riot-api/internal/account"
	"riot-api/internal/match"
	"riot-api/internal/mmr"
	"strconv"
	"strings"
	"time"
)

type requestResult struct {
	Data   json.RawMessage `json:"data"`
	Status int32           `json:"status"`
}

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
func GetAccountByName(apiKey string, params map[string]string, query map[string]any) account.Account {
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v2/account/{name}/{tag}", params, query)

	var account account.Account
	err := json.Unmarshal(result.Data, &account)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(account, result.Status)
	return account
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
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v2/by-puuid/account/{puuid}", params, query)

	var account account.Account
	err := json.Unmarshal(result.Data, &account)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(account, result.Status)
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
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v3/by-puuid/mmr/{affinity}/{platform}/{puuid}", params, nil)

	var mmr mmr.MMR
	err := json.Unmarshal(result.Data, &mmr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(mmr, result.Status)
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
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v4/by-puuid/matches/{affinity}/{platform}/{puuid}", params, query)

	var matches []match.Match
	err := json.Unmarshal(result.Data, &matches)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(matches, result.Status)
}

// GetMatchInfo возвращает данные матча VALORANT по идентификатору матча.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - affinity — регион игрока (e.g., na, eu, ap, kr)
//   - match_id — идентификатор матча
func GetMatchInfo(apiKey string, params map[string]string) {
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v4/match/{affinity}/{match_id}", params, nil)

	var match match.Match
	err := json.Unmarshal(result.Data, &match)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(match, result.Status)
}

// Вспомогательные функции
func getEndpoint(path string, params map[string]string) string {
	endpoint := path

	for k, v := range params {
		templateKey := fmt.Sprintf("{%s}", k)

		endpoint = strings.ReplaceAll(endpoint, templateKey, url.PathEscape(v))
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

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		DisableKeepAlives: true, // Заставляет закрывать TCP-соединение сразу после каждого запроса
	},
}

func baseRequest(apiKey, urlStr string, params map[string]string, query map[string]any) requestResult {
	endpoint := getEndpoint(urlStr, params)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Println("Ошибка создания HTTP-запроса:", err)
		return requestResult{}
	}

	if query != nil {
		urlValues := req.URL.Query()
		setQuery(urlValues, query)
		req.URL.RawQuery = urlValues.Encode()
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Accept", "application/json")

	// Выполняем запрос через наш оптимизированный httpClient
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("КРИТИЧЕСКАЯ ОШИБКА СЕТИ (Do):", err)
		return requestResult{}
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа (Stream broken):", err)
		return requestResult{}
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Сервер вернул ошибку. HTTP Статус: %d, Ответ: %s\n", res.StatusCode, string(bodyBytes))
		return requestResult{Status: int32(res.StatusCode)}
	}

	var result requestResult
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Println("Ошибка декодирования структуры requestResult:", err)
		return requestResult{}
	}

	return result
}
