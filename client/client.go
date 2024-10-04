package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ClientOpts struct {
	Host   string
	Port   int
	UseSsl bool
}

// Возвращает опции подключения к серверу
func NewClientOpts(host string, port int, usessl bool) *ClientOpts {
	return &ClientOpts{
		Host:   host,
		Port:   port,
		UseSsl: usessl,
	}
}

type Client struct {
	ConnectionLine string
	Agent          AgentManager
	User           UserManager
	Group          GroupManager
	Privelege      PrivelegeManager
}

// NewClient создает нового клиента для соединения с микросервисом
func NewClient(opts *ClientOpts) (*Client, error) {
	schema := "http"
	if opts.UseSsl {
		schema = "https"
	}
	connectionLine := fmt.Sprintf("%s://%s:%d", schema, opts.Host, opts.Port)
	return &Client{
		ConnectionLine: connectionLine,
		Agent:          AgentManager{ConnectionLine: connectionLine},
		User:           UserManager{ConnectionLine: connectionLine},
		Group:          GroupManager{ConnectionLine: connectionLine},
		Privelege:      PrivelegeManager{ConnectionLine: connectionLine},
	}, nil
}

// Ping проверяет, отвечает ли сервер. В случае успеха должен вернуть статус 200;
func (c *Client) Ping() error {
	urlRequest := fmt.Sprintf("%s/api/v1/ping", c.ConnectionLine)
	resp, err := http.Get(urlRequest)
	if err != nil {
		return errors.Wrapf(err, "server is not respond at the address %s", c.ConnectionLine)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("server error")
	}
	return nil
}

// //////// AGENT //////////
type AgentManager struct {
	ConnectionLine string
}

// Create создает агента
func (a *AgentManager) Create(agentName, emailCreate string) (Agent, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/agents/%s/who_creates/%s",
		a.ConnectionLine, agentName, emailCreate)
	respRequest, err := http.Post(urlRequest, "application/json", nil)
	if err != nil {
		return Agent{}, err
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp Agent
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return Agent{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return Agent{}, err
		}
		return Agent{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return Agent{}, fmt.Errorf("unexpected error")
	}
}

// Delete удаляет агента
func (a *AgentManager) Delete(agentName, emailDelete string) (ResponseDetail, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/agents/%s/who_deletes/%s",
		a.ConnectionLine, agentName, emailDelete)
	req, err := http.NewRequest("DELETE", urlRequest, nil)
	if err != nil {
		return ResponseDetail{}, err
	}

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return ResponseDetail{}, err
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return ResponseDetail{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return ResponseDetail{}, fmt.Errorf("unexpected error")
	}
}

// GetAll возвращает всех агентов в системе
func (a *AgentManager) GetAll(emailRead string) ([]Agent, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/agents/who_reads/%s", a.ConnectionLine, emailRead)
	respRequest, err := http.Get(urlRequest)
	if err != nil {
		return nil, err
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp []Agent
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error: %s", resp.Error)

	default:
		return nil, fmt.Errorf("unexpected error")
	}
}

// //////// USER //////////
type UserManager struct {
	ConnectionLine string
}

// Create создает пользователя
func (u *UserManager) Create(email, password, firstName, lastName string) (UserWithoutPassword, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/users", u.ConnectionLine)
	body := CreateData{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}
	bodyEncoded, err := json.Marshal(body)
	if err != nil {
		return UserWithoutPassword{}, errors.Wrapf(err, "error while marshalling data")
	}
	respRequest, err := http.Post(urlRequest, "application/json", bytes.NewBuffer(bodyEncoded))
	if err != nil {
		return UserWithoutPassword{}, err
	}

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp UserWithoutPassword
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return UserWithoutPassword{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return UserWithoutPassword{}, err
		}
		return UserWithoutPassword{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return UserWithoutPassword{}, fmt.Errorf("unexpected error")
	}
}

// Get возвращает пользователя
func (u *UserManager) Get(email string) (UserWithoutPassword, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s", u.ConnectionLine, email)
	respRequest, err := http.Get(urlRequest)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp UserWithoutPassword
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return UserWithoutPassword{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return UserWithoutPassword{}, err
		}
		return UserWithoutPassword{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return UserWithoutPassword{}, fmt.Errorf("unexpected error")
	}
}

// Delete удаляет пользователя
func (a *UserManager) Delete(email, emailDelete string) (ResponseDetail, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/who_deletes/%s", a.ConnectionLine, email, emailDelete)
	req, err := http.NewRequest("DELETE", urlRequest, nil)
	if err != nil {
		return ResponseDetail{}, err
	}

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return ResponseDetail{}, err
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return ResponseDetail{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return ResponseDetail{}, fmt.Errorf("unexpected error")
	}
}

// //////// GROUP //////////
type GroupManager struct {
	ConnectionLine string
}

// AddUserToGroup добавляет пользователя в группу
func (g *GroupManager) AddUserToGroup(groupName, email, emailInvite string) (ResponseDetail, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/add_user/%s/who_invites/%s",
		g.ConnectionLine, groupName, email, emailInvite)
	respRequest, err := http.Post(urlRequest, "application/json", nil)
	if err != nil {
		return ResponseDetail{}, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return ResponseDetail{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return ResponseDetail{}, fmt.Errorf("unexpected error")
	}
}

// UserList возвращает группы пользователя
func (g *GroupManager) UserList(email, emailAsk string) ([]Group, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/groups/who_asks/%s", g.ConnectionLine, email, emailAsk)
	respRequest, err := http.Get(urlRequest)
	if err != nil {
		return nil, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var groups []Group
		err = json.NewDecoder(respRequest.Body).Decode(&groups)
		if err != nil {
			return nil, err
		}
		return groups, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error: %s", resp.Error)

	default:
		return nil, fmt.Errorf("unexpected error")
	}
}

// KickOutUser удаляет пользователя из группы
func (g *GroupManager) KickOutUser(groupName, email, emailKick string) (ResponseDetail, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/kick_user/%s/who_kicks/%s",
		g.ConnectionLine, groupName, email, emailKick)
	respRequest, err := http.Post(urlRequest, "application/json", nil)
	if err != nil {
		return ResponseDetail{}, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return ResponseDetail{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return ResponseDetail{}, fmt.Errorf("unexpected error")
	}
}

// MakeBidToCreateGroup создает заявку на создание группы
func (g *GroupManager) MakeBidToCreateGroup(groupName, email string) (Bid, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/who_adds/%s",
		g.ConnectionLine, groupName, email)
	respRequest, err := http.Post(urlRequest, "application/json", nil)
	if err != nil {
		return Bid{}, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var bid Bid
		err = json.NewDecoder(respRequest.Body).Decode(&bid)
		if err != nil {
			return Bid{}, err
		}
		return bid, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return Bid{}, err
		}
		return Bid{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return Bid{}, fmt.Errorf("unexpected error")
	}
}

// ChangeBidStatus меняет статус заявки на создание группы
func (g *GroupManager) ChangeBidStatus(groupName, email, emailChangeStatus, newStatus string) (Bid, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/groups/%s/who_change_status/%s?status=%s",
		g.ConnectionLine, email, groupName, emailChangeStatus, newStatus)
	req, err := http.NewRequest("PUT", urlRequest, nil)
	if err != nil {
		return Bid{}, err
	}

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return Bid{}, err
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp Bid
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return Bid{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return Bid{}, err
		}
		return Bid{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return Bid{}, fmt.Errorf("unexpected error")
	}
}

// ChangeOwner изменяет ответственного в группе
func (g *GroupManager) ChangeOwner(groupName, email, emailWhoChange string) (Group, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/users/%s/who_change_owner/%s",
		g.ConnectionLine, groupName, email, emailWhoChange)
	req, err := http.NewRequest("PUT", urlRequest, nil)
	if err != nil {
		return Group{}, err
	}

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return Group{}, err
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp Group
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return Group{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return Group{}, err
		}
		return Group{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return Group{}, fmt.Errorf("unexpected error")
	}
}

// //////// PRIVELEGE //////////
type PrivelegeManager struct {
	ConnectionLine string
}

// AddAgentToGroup создает связь между агентом и группой
func (p *PrivelegeManager) AddAgentToGroup(groupName, agentName, emailAdd string) (ResponseDetail, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/priveleges/new/agents/%s/who_adds/%s",
		p.ConnectionLine, groupName, agentName, emailAdd)
	respRequest, err := http.Post(urlRequest, "application/json", nil)
	if err != nil {
		return ResponseDetail{}, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return ResponseDetail{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return ResponseDetail{}, fmt.Errorf("unexpected error")
	}
}

// DeleteAgentFromGroup разрывает связь между агентом и группой
func (p *PrivelegeManager) DeleteAgentFromGroup(groupName, agentName, emailDelete string) (ResponseDetail, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/priveleges/delete/agents/%s/who_deletes/%s",
		p.ConnectionLine, groupName, agentName, emailDelete)
	req, err := http.NewRequest("PUT", urlRequest, nil)
	if err != nil {
		return ResponseDetail{}, err
	}

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return ResponseDetail{}, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return resp, nil

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return ResponseDetail{}, err
		}
		return ResponseDetail{}, fmt.Errorf("error: %s", resp.Error)

	default:
		return ResponseDetail{}, fmt.Errorf("unexpected error")
	}
}

// GetGroupAgents возвращает список агентов какойлибо группы
func (p *PrivelegeManager) GetGroupAgents(groupName, emailAsk string) ([]Agent, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/priveleges/who_asks/%s",
		p.ConnectionLine, groupName, emailAsk)
	respRequest, err := http.Get(urlRequest)
	if err != nil {
		return nil, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var agents []Agent
		err = json.NewDecoder(respRequest.Body).Decode(&agents)
		if err != nil {
			return nil, err
		}
		return agents, nil

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error: %s", resp.Error)

	default:
		return nil, fmt.Errorf("unexpected error")
	}
}

// CanUserExecute проверяет, может ли пользователь выполнить процесс на выбранном агенте
func (p *PrivelegeManager) CanUserExecute(email, agentName string) (bool, error) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/check_access/agents/%s",
		p.ConnectionLine, email, agentName)
	respRequest, err := http.Get(urlRequest)
	if err != nil {
		return false, err
	}
	switch respRequest.StatusCode {
	case http.StatusOK:
		var data map[string]bool
		err = json.NewDecoder(respRequest.Body).Decode(&data)
		if err != nil {
			return false, err
		}
		return data["can_execute"], nil

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("error: %s", resp.Error)

	default:
		return false, fmt.Errorf("unexpected error")
	}
}
