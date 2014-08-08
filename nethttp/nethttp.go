package nethttp

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type NetHttp struct {
	Url      string
	Proxy    string
	PostData string
	Header   map[string]string
	Cookie   *cookiejar.Jar
	Timeout  int
}

func NewNetHttp() *NetHttp {
	cj, _ := cookiejar.New(nil)
	header := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3",
		"Connection":      "keep-alive",
		//"Host":            "",
		//"Referer":         "",
		"User-Agent": "Mozilla/5.0 (Linux; U; Android 4.2.2; HTC One Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.90 Mobile Safari/537.36",
	}
	return &NetHttp{"", "", "", header, cj, 60}
}

func (netHttp *NetHttp) NewCookie() {
	netHttp.Cookie, _ = cookiejar.New(nil)
}

func (netHttp *NetHttp) HttpGet() (int64, string, error) {
	ts := time.Now().UnixNano()

	var client *http.Client

	if netHttp.Proxy != "" {
		proxy, err := url.Parse(netHttp.Proxy)
		if err != nil {
			return 0, "", err
		}

		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
				Dial: func(network, addr string) (net.Conn, error) {
					deadline := time.Now().Add(time.Duration(netHttp.Timeout*5) * time.Second)
					c, err := net.DialTimeout(network, addr, time.Second*time.Duration(netHttp.Timeout))
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
			Jar: netHttp.Cookie,
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					deadline := time.Now().Add(time.Duration(netHttp.Timeout*5) * time.Second)
					c, err := net.DialTimeout(network, addr, time.Second*time.Duration(netHttp.Timeout))
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
			Jar: netHttp.Cookie,
		}
	}

	reqest, err := http.NewRequest("GET", netHttp.Url, nil)

	if err != nil {
		return 0, "", err
	}

	for key, value := range netHttp.Header {
		reqest.Header.Add(key, value)
	}

	response, err := client.Do(reqest)

	if err != nil {
		return 0, "", err
	}

	defer response.Body.Close()

	te := time.Now().UnixNano()

	if response.StatusCode == 200 {

		var body string

		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err := gzip.NewReader(response.Body)
			if err != nil {
				return 0, "", err
			}

			buf, err := ioutil.ReadAll(reader)
			if err != nil {
				return 0, "", err
			}
			body = string(buf)
		default:
			bodyByte, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return 0, "", err
			}
			body = string(bodyByte)
		}

		return (te - ts) / 1000000, body, nil
	}

	return (te - ts) / 1000000, "", errors.New(fmt.Sprintf("response.StatusCode:%d", response.StatusCode))
}

func (netHttp *NetHttp) HttpPost() (int64, string, error) {
	ts := time.Now().UnixNano()

	var client *http.Client

	if netHttp.Proxy != "" {
		proxy, err := url.Parse(netHttp.Proxy)
		if err != nil {
			return 0, "", err
		}

		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
				Dial: func(network, addr string) (net.Conn, error) {
					deadline := time.Now().Add(time.Duration(netHttp.Timeout*5) * time.Second)
					c, err := net.DialTimeout(network, addr, time.Second*time.Duration(netHttp.Timeout))
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
			Jar: netHttp.Cookie,
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					deadline := time.Now().Add(time.Duration(netHttp.Timeout*5) * time.Second)
					c, err := net.DialTimeout(network, addr, time.Second*time.Duration(netHttp.Timeout))
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
			Jar: netHttp.Cookie,
		}
	}

	reqest, err := http.NewRequest("POST", netHttp.Url, strings.NewReader(netHttp.PostData))

	if err != nil {
		return 0, "", err
	}

	for key, value := range netHttp.Header {
		reqest.Header.Add(key, value)
	}

	reqest.Header.Add("Content-Length", strconv.Itoa(len(netHttp.PostData)))
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(reqest)

	if err != nil {
		return 0, "", err
	}

	defer response.Body.Close()

	te := time.Now().UnixNano()

	if response.StatusCode == 200 {

		var body string

		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err := gzip.NewReader(response.Body)
			if err != nil {
				if err != nil {
					return 0, "", err
				}
			}
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					return 0, "", err
				}

				if n == 0 {
					break
				}
				body += string(buf)
			}
		default:
			bodyByte, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return 0, "", err
			}
			body = string(bodyByte)
		}

		return (te - ts) / 1000000, body, nil
	}

	return (te - ts) / 1000000, "", errors.New(fmt.Sprintf("response.StatusCode:%d", response.StatusCode))
}
