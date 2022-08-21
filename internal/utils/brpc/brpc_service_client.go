package brpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"
)

const ClientInstanceErrorCode = 10100
const ClientInputMarshalErrorCode = 10100
const ClientBuildRequestErrorCode = 10101
const ClientSendRequestErrorCode = 10102
const ClientInvalidBrpcErrErrorCode = 10103
const ClientResponseBodyErrorCode = 10104
const ClientOutputUnmarshalErrorCode = 10105

type ClientError struct {
	code    int
	message string
}

func (m *ClientError) Code() int {
	return m.code
}

func (m *ClientError) Error() string {
	return m.message
}

func newError(code int, message string) *ClientError {
	return &ClientError{code: code, message: message}
}

type Address struct {
	Service string
	Region  string
	Path    string
}

type Client struct {
	ipListsUpdatingMutex sync.RWMutex
	ipLists              map[string][]string

	innerClient *http.Client
}

func (m *Client) UpdateInstanceIPs(ips map[string][]string) {
	m.ipListsUpdatingMutex.Lock()
	m.ipLists = ips
	m.ipListsUpdatingMutex.Unlock()
}

// GetInstanceIPs
// get the copy of instance ip lists
func (m *Client) GetInstanceIPs() map[string][]string {

	tmpIPLists := make(map[string][]string)

	m.ipListsUpdatingMutex.RLock()
	for key, ipList := range m.ipLists {
		tmpIPList := ipList
		for _, v := range ipList {
			tmpIPList = append(tmpIPList, v)
		}
		tmpIPLists[key] = tmpIPList
	}
	m.ipListsUpdatingMutex.RUnlock()

	return m.ipLists
}

func (m *Client) SelectInstance(service, region string) (string, error) {
	return "", nil
}

func (m *Client) Init() {
	m.innerClient = &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
}

func (m *Client) Do(ctx context.Context, addr *Address, input proto.Message, output proto.Message) *ClientError {

	instance, err := m.SelectInstance(addr.Service, addr.Region)
	if err != nil {
		return newError(ClientInstanceErrorCode, err.Error())
	}
	url := fmt.Sprintf("http://%s/%s", instance, addr.Path)

	ibs, err := proto.Marshal(input)
	if err != nil {
		return newError(ClientInputMarshalErrorCode, err.Error())
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(ibs))
	if err != nil {
		return newError(ClientBuildRequestErrorCode, err.Error())
	}

	response, err := m.innerClient.Do(request)
	if err != nil {
		return newError(ClientSendRequestErrorCode, err.Error())
	}
	defer func() {
		closerErr := response.Body.Close()
		if closerErr != nil {

		}
	}()

	if brpcErrorCode := response.Header.Get("x-bd-error-code"); response.StatusCode == http.StatusServiceUnavailable && brpcErrorCode != "" {
		brpcErrorCode_int, err := strconv.Atoi(brpcErrorCode)
		if err != nil {
			return newError(ClientInvalidBrpcErrErrorCode, fmt.Sprintf("brpc error code is %s", brpcErrorCode))
		}
		return newError(brpcErrorCode_int, fmt.Sprintf("brpc error code is %s", brpcErrorCode))
	}

	if response.StatusCode != http.StatusOK {
		return newError(response.StatusCode, fmt.Sprintf("response status code is %d", response.StatusCode))
	}

	obs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return newError(ClientResponseBodyErrorCode, err.Error())
	}

	err = proto.Unmarshal(obs, output)
	if err != nil {
		return newError(ClientOutputUnmarshalErrorCode, err.Error())
	}

	return nil
}
