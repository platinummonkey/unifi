package unifi

import (
	"crypto/tls"
	"crypto/x509"
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
	baseURLStr string
	baseURL    *url.URL
	certConfig *CertificationConfig

	HTTPClient   *http.Client
	RetryTimeout time.Duration

	authCookies        []*http.Cookie
	longRunningSession bool
}

// CertificationConfig overrides the default HTTP client behavior with certificates.
// In most home default installations, you're probably using the default self-signed cert that came with the device.
// This setting allows you to override with your cert validation so that you are not running without certificate checks.
// The options provided are mutually exclusive to each other.
// 1. If DisableCertCheck is true, then all certificate checkin is disabled.
// 2. If PEMCert is provided it will use that certificate.
// 3. If Certificates are provided, then it will use those provided certificates.
// 4. The default behavior if nothing is configured or a `nil` CertificationConfig is passed, then the default
//    go http certificate checks are used.
type CertificationConfig struct {
	DisableCertCheck bool                // set to true to disable all certificate checks
	PEMCert          string              // path to custom cert for self-signed certs
	Certificates     []*x509.Certificate // custom certificates to add
}

// NewClient will create a new UniFi http(s) client.
func NewClient(baseURL string, certConfig *CertificationConfig, timeout time.Duration) (*Client, error) {
	httpClient := http.DefaultClient
	if certConfig != nil {
		defaultTransport := http.DefaultTransport.(*http.Transport)
		var tlsConfig *tls.Config
		if certConfig.DisableCertCheck {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else if certConfig.PEMCert != "" {
			cert, err := ioutil.ReadFile(certConfig.PEMCert)
			if err != nil {
				return nil, err
			}
			certPool := x509.NewCertPool()
			certPool.AppendCertsFromPEM(cert)
			tlsConfig = &tls.Config{
				RootCAs: certPool,
			}
		} else if certConfig.Certificates != nil && len(certConfig.Certificates) > 0 {
			certPool := x509.NewCertPool()
			for _, cert := range certConfig.Certificates {
				if cert != nil {
					certPool.AddCert(cert)
				}
			}
			tlsConfig = &tls.Config{
				RootCAs: certPool,
			}
		}

		tr := &http.Transport{
			Proxy:                 defaultTransport.Proxy,
			DialContext:           defaultTransport.DialContext,
			MaxIdleConns:          defaultTransport.MaxIdleConns,
			IdleConnTimeout:       defaultTransport.IdleConnTimeout,
			ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
			TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
			TLSClientConfig:       tlsConfig,
		}
		httpClient = &http.Client{Transport: tr}
	}
	httpClient.Timeout = timeout

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURLStr:   baseURL,
		baseURL:      u,
		certConfig:   certConfig,
		HTTPClient:   httpClient,
		RetryTimeout: timeout,
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

// SetHeaders will set the client headers for auth along with additional pre-defined API headers.
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

// WithPathAndQueryParams will return a normalized url with the baseURL included.
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

// ResponseCodeTrait defines the interface that returns the response code
type ResponseCodeTrait interface {
	GetResponseCode() ResponseCode
}

// ResponseMessageTrait defines an interface that returns the response message.
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if !rv.IsNil() {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, ErrInvalidResponseBody.Error())
		}

		err = json.Unmarshal(body, ret)
		if err != nil {
			return errors.Wrap(err, ErrJSONDecode.Error())
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
