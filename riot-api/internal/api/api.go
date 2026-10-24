package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
func GetAccountByName(apiKey string, params map[string]string, query map[string]any) (json.RawMessage, error) {
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v2/account/{name}/{tag}", params, query)

	if result.Status != 200 {
		return nil, fmt.Errorf("запрос не удался")
	}

	return result.Data, nil
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
func GetAccount(apiKey string, params map[string]string, query map[string]any) (json.RawMessage, error) {
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v2/by-puuid/account/{puuid}", params, query)

	if result.Status != 200 {
		return nil, fmt.Errorf("запрос не удался")
	}

	return result.Data, nil
}

// GetPlayerMMR возвращает рейтинг пользователя VALORANT по идентификатору пользователя.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - affinity — регион игрока (e.g., na, eu, ap, kr)
//   - platform — платформа игрока (pc, console)
//   - puuid — идентификатор игрока
func GetPlayerMMR(apiKey string, params map[string]string) (json.RawMessage, error) {
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v3/by-puuid/mmr/{affinity}/{platform}/{puuid}", params, nil)

	if result.Status != 200 {
		return nil, fmt.Errorf("запрос не удался")
	}

	return result.Data, nil
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
func GetPlayerMatches(apiKey string, params map[string]string, query map[string]any) (json.RawMessage, error) {
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v4/by-puuid/matches/{affinity}/{platform}/{puuid}", params, query)

	if result.Status != 200 {
		return nil, fmt.Errorf("запрос не удался")
	}

	return result.Data, nil
}

// GetMatchInfo возвращает данные матча VALORANT по идентификатору матча.
//
// apiKey — API-ключ HenrikDev.
//
// params содержит path-параметры:
//   - affinity — регион игрока (e.g., na, eu, ap, kr)
//   - match_id — идентификатор матча
func GetMatchInfo(apiKey string, params map[string]string) (json.RawMessage, error) {
	result := baseRequest(apiKey, "https://api.henrikdev.xyz/valorant/v4/match/{affinity}/{match_id}", params, nil)

	if result.Status != 200 {
		return nil, fmt.Errorf("запрос не удался")
	}

	return result.Data, nil
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

// retryableStatusCodes содержит коды, при которых имеет смысл сделать повторный запрос
var retryableStatusCodes = map[int]bool{
	http.StatusTooManyRequests:     true, // 429
	http.StatusInternalServerError: true, // 500
	http.StatusBadGateway:          true, // 502
	http.StatusServiceUnavailable:  true, // 503
	http.StatusGatewayTimeout:      true, // 504
}

func baseRequest(apiKey, urlStr string, params map[string]string, query map[string]any) requestResult {
	endpoint := getEndpoint(urlStr, params)

	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
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

		res, err := httpClient.Do(req)
		if err != nil {
			// Сетевая ошибка (например, таймаут или разрыв соединения)
			if attempt < maxRetries {
				delay := baseDelay * time.Duration(1<<uint(attempt)) // Экспоненциальная задержка: 1s, 2s, 4s
				fmt.Printf("Сетевая ошибка (попытка %d/%d). Повтор через %v: %v\n", attempt+1, maxRetries, delay, err)
				time.Sleep(delay)
				continue
			}
			fmt.Println("КРИТИЧЕСКАЯ ОШИБКА СЕТИ (исчерпаны попытки):", err)
			return requestResult{}
		}

		bodyBytes, _ := io.ReadAll(res.Body)
		res.Body.Close() // Закрываем тело сразу после чтения

		if retryableStatusCodes[res.StatusCode] && attempt < maxRetries {
			// Определяем время задержки
			delay := baseDelay * time.Duration(1<<uint(attempt))

			// Проверяем заголовок Retry-After (может быть в секундах или в виде даты)
			if retryAfter := res.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					delay = time.Duration(seconds) * time.Second
				} else if t, err := http.ParseTime(retryAfter); err == nil {
					delay = time.Until(t)
					if delay < 0 {
						delay = baseDelay // Если время уже прошло, используем базовую задержку
					}
				}
			}

			fmt.Printf("Получен код %d. Повторный запрос через %v (попытка %d/%d)\n", res.StatusCode, delay, attempt+1, maxRetries)
			time.Sleep(delay)
			continue
		}

		// Если код не из списка retryable или попытки исчерпаны
		if res.StatusCode != http.StatusOK {
			fmt.Printf("Сервер вернул ошибку. HTTP Статус: %d, Ответ: %s\n", res.StatusCode, string(bodyBytes))
			return requestResult{Status: int32(res.StatusCode)}
		}

		var result requestResult
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			fmt.Println("Ошибка декодирования структуры requestResult:", err)
			return requestResult{}
		}

		return result
	}

	return requestResult{}
}
