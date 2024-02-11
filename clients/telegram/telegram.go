package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"go_rest_app/main/lib/e"
	"go_rest_app/main/lib/settings"
)

type Client struct {
	host string
	basePath string
	channelID int
	client http.Client
}

const (
	getUpdatesMethod = "getUpdates"
	sendMessageMethod = "sendMessage"
	bot_host = "api.telegram.org"
)

func New(s *settings.Telegram) *Client {
	return &Client {
		host: bot_host,
		basePath: newBasePath(s.Token),
		channelID: s.ChatID,
		client: http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offsest int, limit int) (updates []Update, err error ) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offsest))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.makeRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(text string) error {	
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(c.channelID))
	q.Add("text", text)

	_, err := c.makeRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) makeRequest(method string, query url.Values) (data []byte, err error) {
	defer func () { err = e.WrapIfErr("can't make request",err) }()

	u := url.URL {
		Scheme: "https",
		Host: c.host,
		Path: path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func () { _= resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}