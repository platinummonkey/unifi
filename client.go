package unifi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

// Client is the object that handles talking to the Unifi Controller API. This maintains
// state information for a particular application connection.
type Client struct {
	baseURLStr              string
	baseURL                 *url.URL
	disableCertificateCheck bool

	// The HTTP Client that is used to make requests
	HttpClient   *http.Client
	RetryTimeout time.Duration

	authCookies        []*http.Cookie
	longRunningSession bool
}

func NewClient(baseURL string, disableCertificateCheck bool) (*Client, error) {
	httpClient := http.DefaultClient
	if disableCertificateCheck {
		defaultTransport := http.DefaultTransport.(*http.Transport)

		tr := &http.Transport{
			Proxy:                 defaultTransport.Proxy,
			DialContext:           defaultTransport.DialContext,
			MaxIdleConns:          defaultTransport.MaxIdleConns,
			IdleConnTimeout:       defaultTransport.IdleConnTimeout,
			ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
			TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: tr}
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURLStr:              baseURL,
		baseURL:                 u,
		disableCertificateCheck: disableCertificateCheck,
		HttpClient:              httpClient,
		RetryTimeout:            time.Second * 30,
	}, nil
}

// SetBaseURL changes the value of baseURL.
func (c *Client) SetBaseURL(baseURL string) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	c.baseURLStr = baseURL
	c.baseURL = u
	return nil
}

// GetBaseURL returns the baseURL.
func (c *Client) GetBaseURL() string {
	return c.baseURLStr
}

func (c *Client) SetHeaders(r *http.Request) {
	r.Header.Set("Content-Type", ContentTypeHeader)
	r.Header.Set("Accept", ContentTypeHeader)
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Accept-Charset", "utf-8")
	r.Header.Set("User-Agent", UserAgentHeader)
	if c.authCookies != nil {
		for _, cookie := range c.authCookies {
			r.AddCookie(cookie)
		}
	}
}

func (c *Client) WithPathAndQueryParams(extPath string, queryParamsPairs ...string) *url.URL {
	newPath := path.Join(c.baseURL.Path, extPath)
	u := &url.URL{
		Scheme: c.baseURL.Scheme,
		Host:   c.baseURL.Host,
		Path:   newPath,
	}
	q := u.Query()
	if len(queryParamsPairs)%2 == 0 {
		for i := 0; i < len(queryParamsPairs); i = i + 2 {
			q.Set(queryParamsPairs[i], queryParamsPairs[i+1])
		}
	}
	u.RawQuery = q.Encode()
	return u
}

type ResponseCodeTrait interface {
	GetResponseCode() ResponseCode
}

type ResponseMessageTrait interface {
	GetResponseMessage() string
}

func (c *Client) doRequest(method string, extPath string, sendBody io.Reader, ret interface{}, queryParamsPairs ...string) error {
	u := c.WithPathAndQueryParams(extPath)

	rv := reflect.ValueOf(ret)
	if !rv.IsNil() && rv.Kind() != reflect.Ptr {
		return fmt.Errorf("non nil-response handlers should be a pointer: kind:%v nil:%t", rv.Kind(), rv.IsNil())
	}

	req, err := http.NewRequest(method, u.String(), sendBody)
	if err != nil {
		return err
	}
	c.SetHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if !rv.IsNil() {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, InvalidResponseBody.Error())
		}

		err = json.Unmarshal(body, ret)
		if err != nil {
			return errors.Wrap(err, JSONDecodeError.Error())
		}
		if rv.Kind() == reflect.Struct {
			metaField := rv.FieldByName("Meta")

			if metaField.IsValid() {
				if retRespCodeTrait, ok := metaField.Interface().(ResponseCodeTrait); ok {
					rc := retRespCodeTrait.GetResponseCode()
					if !rc.Equal(ResponseCodeOK) {
						if retRespCodeMsgTrait, ok := metaField.Interface().(ResponseMessageTrait); ok {
							msg := retRespCodeMsgTrait.GetResponseMessage()
							if msg != "" {
								return fmt.Errorf(msg)
							}
						}
						return fmt.Errorf("non-ok status code: %v - %v", rc, ret)
					}
				}
			}
		}
	}

	return nil
}

func (c *Client) doSiteRequest(method string, site string, extPath string, sendBody io.Reader, ret interface{}, queryParamsPairs ...string) error {
	return c.doRequest(method, fmt.Sprintf("/api/s/%s/%s", site, extPath), sendBody, ret, queryParamsPairs...)
}
