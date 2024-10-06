package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ClientOpts struct {
	Host   string
	Port   int
	UseSsl bool
}

const (
	XRealIP   = "X-Real-IP"
	UserAgent = "User-Agent"
)

var (
	ErrInternal = errors.New("internal server error, please try again later")
)

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
	}
}

func (c *Client) CheckConnection() {
	logger := zap.Must(zap.NewProduction())
	i := 0
	for ; i < 3; i++ {
		if err := c.Ping(); err == nil {
			break
		}
		logger.Warn("error while connecting to microservice 'archive'")
		if i < 2 {
			time.Sleep(2 * time.Second)
		}
	}
	if i == 3 {
		logger.Fatal("failed to connect to microservice 'archive'")
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

func (c *Client) GetArchive(meta *RequestMeta) ([]Record, *RequestStatus) {
	urlRequest := fmt.Sprintf("%s/api/v1/archives", c.ConnectionLine)
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
		var resp []Record
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
