package awsauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func serviceAndRegion(host string) (string, string) {
	var region, service string
	parts := strings.Split(host, ".")

	service = parts[0]

	if len(parts) >= 4 {
		if parts[1] == "s3" {
			region = parts[0]
			service = parts[1]
		} else {
			region = parts[1]
		}
	} else {
		if strings.HasPrefix(parts[0], "s3-") {
			service = parts[0][:2]
			region = parts[0][3:]
		} else {
			region = "us-east-1" // default. http://docs.aws.amazon.com/general/latest/gr/rande.html
		}
	}

	return service, region
}

func checkKeys() {
	if Keys == nil {
		Keys = &Credentials{
			os.Getenv(envAccessKeyID),
			os.Getenv(envSecretAccessKey),
		}
	}
}

func augmentRequestQuery(req *http.Request, values url.Values) *http.Request {
	for key, arr := range req.URL.Query() {
		for _, val := range arr {
			values.Set(key, val)
		}
	}

	req.URL.RawQuery = values.Encode()

	return req
}

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func hmacSHA1(key []byte, content string) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func hashSHA256(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func readAndReplaceBody(req *http.Request) string {
	rawPayload, _ := ioutil.ReadAll(req.Body)
	payload := string(rawPayload)
	req.Body = ioutil.NopCloser(strings.NewReader(payload))
	return payload
}

func concat(delim string, str ...string) string {
	return strings.Join(str, delim)
}

var now = func() time.Time {
	return time.Now().UTC()
}