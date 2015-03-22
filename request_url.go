package pusher

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"
)

const auth_version = "1.0"
const domain = "http://api.pusherapp.com"

func createBodyMD5(body []byte) string {
	_body_md5 := md5.New()
	_body_md5.Write([]byte(body))
	return hex.EncodeToString(_body_md5.Sum(nil))
}

func unsigned_params(key, timestamp string, body []byte, additional_queries map[string]string) url.Values {
	params := url.Values{
		"auth_key":       {key},
		"auth_timestamp": {timestamp},
		"auth_version":   {auth_version},
	}

	if body != nil {
		params.Add("body_md5", createBodyMD5(body))
	}

	if additional_queries != nil {
		for key, values := range additional_queries {
			params.Add(key, values)
		}
	}

	return params

}

func unescape_url(_url url.Values) string {
	unesc, _ := url.QueryUnescape(_url.Encode())
	return unesc
}

func createRequestUrl(method, path, key, secret, timestamp string, body []byte, additional_queries map[string]string) string {
	params := unsigned_params(key, timestamp, body, additional_queries)

	string_to_sign := strings.Join([]string{method, path, unescape_url(params)}, "\n")

	auth_signature := HMACSignature(string_to_sign, secret)

	params.Add("auth_signature", auth_signature)

	endpoint, _ := url.Parse(domain + path)

	endpoint.RawQuery = unescape_url(params)

	return endpoint.String()
}
