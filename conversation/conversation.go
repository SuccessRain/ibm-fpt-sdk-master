package conversation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultMajorVersion = "v1"
const defaultMinorVersion = "2017-05-05"
//const defaultUrl = "https://gateway.watsonplatform.net/conversation-experimental/api"
const defaultUrl = "https://gateway.watsonplatform.net/conversation/api/v1"

// Connects to instance of Watson Conversation service
func NewClient(cfg watson.Config) (Client, error) {
	ci := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "conversation"
	}
	if len(cfg.Credentials.Url) == 0 {
		cfg.Credentials.Url = defaultUrl
	}
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	ci.watsonClient = client
	return ci, nil
}

type Intent struct {
	Intent     string  `json:"intent,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

type IntentExample struct {
	Text     string          `json:"text,omitempty"`
	Entities []EntityExample `json:"entities,omitempty"`
}

type EntityExample struct {
	Entity   string `json:"entity,omitempty"`
	Value    string `json:"value,omitempty"`
	Location []int  `json:"location,omitempty"`
}

type Message struct {
	Input   MessageInput           `json:"input,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

type MessageInput struct {
	Text string `json:"text,omitempty"`
}

type MessageOutput struct {
	LogMessages []interface{} `json:"log_messages,omitempty"`
	Text        []string      `json:"text,omitempty"`
	HitNodes    []string      `json:"hit_nodes,omitempty"`
}

type MessageResponse struct {
	Input    MessageInput           `json:"input,omitempty"`
	Intents  []Intent               `json:"intents,omitempty"`
	Entities []EntityExample        `json:"entities,omitempty"`
	Output   MessageOutput          `json:"output,omitempty"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// Calls 'GET /v1/workspaces/{workspace_id}/message' to retrieve response from conversation utterance
func (c Client) Message(workspace_id string, text string) (MessageResponse, error) {
	q := url.Values{}
	q.Set("version", defaultMinorVersion)

	message := &Message{Input: MessageInput{Text: text}}
	message_json, err := json.Marshal(message)

	//fmt.Print("Marshal:\t"); fmt.Println(text, message_json)

	if err != nil {
		return MessageResponse{}, err
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	//body, err := c.watsonClient.MakeRequest("POST", c.version+"/workspaces/"+workspace_id+"/message?"+q.Encode(), bytes.NewReader(message_json), headers)
	body, err := c.watsonClient.MakeRequest("POST", "/workspaces/"+workspace_id+"/message?"+q.Encode(), bytes.NewReader(message_json), headers)
	//body, err := c.watsonClient.MakeRequest("POST", "/workspaces/"+workspace_id+"?"+q.Encode(), bytes.NewReader(message_json), headers)
	//body, err := c.watsonClient.MakeRequest("POST", "/workspaces?"+q.Encode(), bytes.NewReader(message_json), headers)

	//fmt.Print("MakeRequest:\t"); fmt.Println(err, body)

	if err != nil {
		return MessageResponse{}, err
	}
	var response MessageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return MessageResponse{}, err
	}

	//fmt.Print("MessageResponse:\t"); fmt.Println(err, response)

	return response, nil
}
