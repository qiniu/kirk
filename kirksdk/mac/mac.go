package mac

import (
	"encoding/base64"
	"net/http"

	"qiniupkg.com/api.v7/conf"
)

type Mac struct {
	AccessKey string
	SecretKey []byte
}

func (m *Mac) SignRequest(req *http.Request) (err error) {

	sign, err := signRequest(m.SecretKey, req)
	if err != nil {
		return
	}

	auth := "Qiniu " + m.AccessKey + ":" + base64.URLEncoding.EncodeToString(sign)
	req.Header.Set("Authorization", auth)
	return
}

type Transport struct {
	mac       Mac
	Transport http.RoundTripper
}

func New(accessKey, secretKey string) *Mac {

	if accessKey == "" {
		accessKey = conf.ACCESS_KEY
		secretKey = conf.SECRET_KEY
	}

	return &Mac{accessKey, []byte(secretKey)}
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	err = t.mac.SignRequest(req)
	if err != nil {
		return
	}

	return t.Transport.RoundTrip(req)
}

func NewTransport(mac *Mac, transport http.RoundTripper) *Transport {

	if transport == nil {
		transport = http.DefaultTransport
	}
	t := &Transport{Transport: transport}
	if mac == nil {
		t.mac.AccessKey = conf.ACCESS_KEY
		t.mac.SecretKey = []byte(conf.SECRET_KEY)
	} else {
		t.mac = *mac
	}
	return t
}

func NewClient(mac *Mac, transport http.RoundTripper) *http.Client {

	t := NewTransport(mac, transport)
	return &http.Client{Transport: t}
}
