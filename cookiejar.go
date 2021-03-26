package areq

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type CookieJar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewCookieJar() *CookieJar {
	jar := new(CookieJar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func NewCookieJarFromFile(filename string) (*CookieJar, error) {
	j := new(CookieJar)
	j.cookies = make(map[string][]*http.Cookie)
	if err := j.LoadFile(filename); err != nil {
		return nil, err
	}
	return j, nil
}

func (j *CookieJar) Load(c []byte) error {
	tmp := make(map[string][]http.Cookie)
	if err := json.Unmarshal(c, &tmp); err != nil {
		return err
	}
	// convert tmp to cookie
	for k, v := range tmp {
		var tmpCookies []*http.Cookie
		for _, v2 := range v {
			tmpCookies = append(tmpCookies, &v2)
		}
		j.cookies[k] = tmpCookies
	}
	return nil
}

func (j *CookieJar) Save() ([]byte, error) {
	return json.Marshal(j)
}

func (j *CookieJar) LoadFile(filename string) error {
	r, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return j.Load(r)
}

func (j *CookieJar) SaveFile(filename string) error {
	b, err := j.Save()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0666)
}

func (j *CookieJar) AllCookies() map[string][]*http.Cookie {
	return j.cookies
}

func (j *CookieJar) List() []string {
	var urls []string
	for k, _ := range j.cookies {
		urls = append(urls, k)
	}
	return urls
}

func (j *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.lk.Lock()
	j.cookies[u.Host] = cookies
	j.lk.Unlock()
}

func (j *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	return j.cookies[u.Host]
}
