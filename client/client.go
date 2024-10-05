package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	XRealIP   = "X-Real-IP"
	UserAgent = "User-Agent"
)

var (
	ErrInternal = errors.New("internal server error, please try again later")
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
func NewClient(opts *ClientOpts) *Client {
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
	}
}

func (c *Client) CheckConnection() {
	logger := zap.Must(zap.NewProduction())
	i := 0
	for ; i < 3; i++ {
		if err := c.Ping(); err == nil {
			break
		}
		logger.Warn("error while connecting to microservice 'privelege'")
		if i < 2 {
			time.Sleep(2 * time.Second)
		}
	}
	if i == 3 {
		logger.Fatal("failed to connect to microservice 'privelege'")
	}
}

// Ping проверяет, отвечает ли сервер. В случае успеха должен вернуть статус 200;
func (c *Client) Ping() error {
	urlRequest := fmt.Sprintf("%s/api/v1/ping", c.ConnectionLine)
	resp, err := http.Get(urlRequest)
	if err != nil {
		return fmt.Errorf("server is not respond at the address %s", c.ConnectionLine)
	}
	if resp.StatusCode != http.StatusOK {
		return ErrInternal
	}
	return nil
}

// //////// AGENT //////////
type AgentManager struct {
	ConnectionLine string
}

// Create создает агента
func (a *AgentManager) Create(agentName, emailCreate string, meta *RequestMeta) (*Agent, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/agents/%s/who_creates/%s",
		a.ConnectionLine, agentName, emailCreate)
	req, err := http.NewRequest("POST", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}

	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp Agent
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)
	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// Delete удаляет агента
func (a *AgentManager) Delete(agentName, emailDelete string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/agents/%s/who_deletes/%s",
		a.ConnectionLine, agentName, emailDelete)
	req, err := http.NewRequest("DELETE", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}

	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// GetAll возвращает всех агентов в системе
func (a *AgentManager) GetAll(emailRead string, meta *RequestMeta) ([]Agent, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/agents/who_reads/%s", a.ConnectionLine, emailRead)
	req, err := http.NewRequest("GET", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp []Agent
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// //////// GROUP //////////
type GroupManager struct {
	ConnectionLine string
}

// AddUserToGroup добавляет пользователя в группу
func (g *GroupManager) AddUserToGroup(groupName, email, emailInvite string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/add_user/%s/who_invites/%s",
		g.ConnectionLine, groupName, email, emailInvite)
	req, err := http.NewRequest("POST", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// UserList возвращает группы пользователя
func (g *GroupManager) UserList(email, emailAsk string, meta *RequestMeta) ([]Group, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/groups/who_asks/%s", g.ConnectionLine, email, emailAsk)
	req, err := http.NewRequest("GET", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp []Group
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// KickOutUser удаляет пользователя из группы
func (g *GroupManager) KickOutUser(groupName, email, emailKick string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/kick_user/%s/who_kicks/%s",
		g.ConnectionLine, groupName, email, emailKick)
	req, err := http.NewRequest("POST", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// MakeBidToCreateGroup создает заявку на создание группы
func (g *GroupManager) MakeBidToCreateGroup(groupName, email string, meta *RequestMeta) (*Bid, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/who_adds/%s",
		g.ConnectionLine, groupName, email)
	req, err := http.NewRequest("POST", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp Bid
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// ChangeBidStatus меняет статус заявки на создание группы
func (g *GroupManager) ChangeBidStatus(groupName, email, emailChangeStatus, newStatus string, meta *RequestMeta) (*Bid, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/groups/%s/who_change_status/%s?status=%s",
		g.ConnectionLine, email, groupName, emailChangeStatus, newStatus)
	req, err := http.NewRequest("PUT", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp Bid
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// ChangeOwner изменяет ответственного в группе
func (g *GroupManager) ChangeOwner(groupName, email, emailWhoChange string, meta *RequestMeta) (*Group, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/users/%s/who_change_owner/%s",
		g.ConnectionLine, groupName, email, emailWhoChange)
	req, err := http.NewRequest("PUT", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp Group
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// //////// USER //////////
type UserManager struct {
	ConnectionLine string
}

// Create создает пользователя
func (u *UserManager) Create(body io.ReadCloser, meta *RequestMeta) (*UserWithoutPassword, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users", u.ConnectionLine)
	req, err := http.NewRequest("POST", urlRequest, body)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp UserWithoutPassword
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// Get возвращает пользователя
func (u *UserManager) Get(email string, meta *RequestMeta) (*UserWithoutPassword, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s", u.ConnectionLine, email)
	req, err := http.NewRequest("GET", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp UserWithoutPassword
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// Delete удаляет пользователя
func (a *UserManager) Delete(email, emailDelete string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/who_deletes/%s", a.ConnectionLine, email, emailDelete)
	req, err := http.NewRequest("DELETE", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// //////// PRIVELEGE //////////
type PrivelegeManager struct {
	ConnectionLine string
}

// AddAgentToGroup создает связь между агентом и группой
func (p *PrivelegeManager) AddAgentToGroup(groupName, agentName, emailAdd string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/priveleges/new/agents/%s/who_adds/%s",
		p.ConnectionLine, groupName, agentName, emailAdd)
	req, err := http.NewRequest("POST", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// DeleteAgentFromGroup разрывает связь между агентом и группой
func (p *PrivelegeManager) DeleteAgentFromGroup(groupName, agentName, emailDelete string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/priveleges/delete/agents/%s/who_deletes/%s",
		p.ConnectionLine, groupName, agentName, emailDelete)
	req, err := http.NewRequest("DELETE", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// GetGroupAgents возвращает список агентов какойлибо группы
func (p *PrivelegeManager) GetGroupAgents(groupName, emailAsk string, meta *RequestMeta) ([]Agent, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/groups/%s/priveleges/who_asks/%s",
		p.ConnectionLine, groupName, emailAsk)

	req, err := http.NewRequest("GET", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp []Agent
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// AddAgentToUser создает связь между агентом и пользователем
func (p *PrivelegeManager) AddAgentToUser(email, agentName, emailAdd string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/priveleges/new/agents/%s/who_adds/%s",
		p.ConnectionLine, email, agentName, emailAdd)
	req, err := http.NewRequest("POST", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// DeleteAgentFromUser разрывает связь между агентом и пользователем
func (p *PrivelegeManager) DeleteAgentFromUser(email, agentName, emailDelete string, meta *RequestMeta) (*ResponseDetail, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/priveleges/delete/agents/%s/who_deletes/%s",
		p.ConnectionLine, email, agentName, emailDelete)
	req, err := http.NewRequest("DELETE", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp ResponseDetail
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return &resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// GetGroupAgents возвращает список агентов какойлибо группы
func (p *PrivelegeManager) GetUserAgents(email, emailAsk string, meta *RequestMeta) ([]Agent, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/priveleges/who_asks/%s",
		p.ConnectionLine, email, emailAsk)
	req, err := http.NewRequest("GET", urlRequest, nil)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var resp []Agent
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return resp, newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return nil, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return nil, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}

// CanUserExecute проверяет, может ли пользователь выполнить процесс на выбранном агенте
func (p *PrivelegeManager) CanUserExecute(email, agentName string, meta *RequestMeta) (bool, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/users/%s/check_access/agents/%s",
		p.ConnectionLine, email, agentName)
	req, err := http.NewRequest("GET", urlRequest, nil)
	if err != nil {
		return false, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	req.Header.Set(XRealIP, meta.RealIp)
	req.Header.Set(UserAgent, meta.UserAgent)

	client := &http.Client{}
	respRequest, err := client.Do(req)
	if err != nil {
		return false, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
	defer respRequest.Body.Close()

	switch respRequest.StatusCode {
	case http.StatusOK:
		var data map[string]bool
		err = json.NewDecoder(respRequest.Body).Decode(&data)
		if err != nil {
			return false, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return data["can_execute"], newRequestStatus(nil, respRequest.StatusCode)

	case http.StatusBadRequest, http.StatusInternalServerError:
		var resp ResponseError
		err = json.NewDecoder(respRequest.Body).Decode(&resp)
		if err != nil {
			return false, newRequestStatus(ErrInternal, http.StatusInternalServerError)
		}
		return false, newRequestStatus(errors.New(resp.Error), respRequest.StatusCode)

	default:
		return false, newRequestStatus(ErrInternal, http.StatusInternalServerError)
	}
}
