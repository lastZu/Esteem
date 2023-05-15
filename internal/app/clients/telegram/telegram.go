package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/lastZu/Esteem/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (client *Client) Updates(offset int, limit int) ([]Update, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := client.doRequest(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	var result UpdatesResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (client *Client) SendMessage(chatID int, text string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)

	_, err := client.doRequest(sendMessageMethod, query)
	if err != nil {
		return e.Wrap("can't send message", err)
	}
	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (client *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	const errMsg = "can't do request"
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	requestUrl := url.URL{
		Scheme: "https",
		Host:   client.host,
		Path:   path.Join(client.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = query.Encode()

	response, err := client.client.Do(request)

	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
