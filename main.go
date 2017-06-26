package main

import (
	"os"
	"fmt"
	"path"
	"io/ioutil"
	"github.com/liuyibao/api-proxy/route"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"crypto/sha1"
)

var AppPath string

var proxyCache map[string][]byte

func main()  {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	AppPath = path.Dir(ex)

	b, err := ioutil.ReadFile(AppPath + string(os.PathSeparator) + "config/route.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	routeMap, err := route.New(b)
	if err != nil {
		log.Fatalln(err)
	}

	go cleanProxyCache()

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request){
		code := 404
		body := []byte(request.URL.Path)
		for _, route := range routeMap {
			if matched, err := regexp.MatchString("^" + route.Path + "$", request.URL.Path); matched && strings.ToUpper(route.Method) == request.Method {
				code = 200
				body = internalRequest(request)
				break
			} else {
				log.Println(err)
			}
		}
		response.WriteHeader(code)
		response.Write(body)
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
	fmt.Println(routeMap)
}

func cleanProxyCache () {
	for {
		time.Sleep(10000 * time.Millisecond)
		fmt.Println(proxyCache)
		proxyCache = make(map[string][]byte)
		fmt.Println(proxyCache)
	}
}

func internalRequest (request *http.Request) []byte {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(request.RequestURI))
	hash := string(Sha1Inst.Sum([]byte{}))
	cache, exist := proxyCache[hash]
	if exist && len(cache) > 0 {
		return cache
	} else {
		resp, err := http.Get("http://api.ffan.com" + request.RequestURI)
		if err != nil {
			log.Println("内部代理错误")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		proxyCache[hash] = body;
		return body;
	}
}
