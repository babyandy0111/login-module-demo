package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type APIResponse struct {
	Data  []interface{} `json:"data"`
	Total int           `json:"total"`
}

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {

		token := "123"
		f, err := os.Open("test.html")
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		tmpl := getTmpl(f, token)
		escaped := html.UnescapeString(tmpl)
		// log.Println("escaped", escaped)
		newHtml := getNewHtml(escaped)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(newHtml)
	})
	app.Listen(":8080")
}

func getTmpl(html *os.File, token string) string {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	dom.Find("range").Each(func(i int, s *goquery.Selection) {
		aliasName := s.AttrOr("data-gjs-aliasname", "")
		endpoint := s.AttrOr("data-gjs-endpoint", "")
		request := s.AttrOr("data-gjs-request", "")
		action := s.AttrOr("data-gjs-action", "")

		content, err := s.Html()
		if err != nil {
			log.Println(err)
			return
		}
		s.ReplaceWithHtml(rangeTagToTmpl(content, aliasName, endpoint, request, action, token))
	})

	dom.Find("for-var").Each(func(i int, s *goquery.Selection) {
		aliasName := s.AttrOr("data-gjs-varname", "")
		content, err := s.Html()
		if err != nil {
			log.Println(err)
			return
		}
		s.ReplaceWithHtml(forVarTagToTmpl(content, aliasName))
	})

	dom.Find("if").Each(func(i int, s *goquery.Selection) {
		logic := s.AttrOr("data-gjs-logic", "")
		content, err := s.Html()
		if err != nil {
			log.Println(err)
			return
		}
		s.ReplaceWithHtml(ifTagToTmpl(content, logic))
	})

	dom.Find("fetchdata").Each(func(i int, s *goquery.Selection) {
		aliasName := s.AttrOr("data-gjs-aliasname", "")
		endpoint := s.AttrOr("data-gjs-endpoint", "")
		request := s.AttrOr("data-gjs-request", "")
		action := s.AttrOr("data-gjs-action", "GET")
		s.ReplaceWithHtml(fetchDataTagToTmpl(aliasName, endpoint, request, action, token))
	})

	tmpl, err := dom.Html()
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return tmpl
}

func json2query(jsonString string) string {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	queryString := ""
	for key, value := range jsonData {
		// fmt.Println("Key:", key, "Value:", value)
		queryString = queryString + fmt.Sprintf("%s=%v&", key, value)
	}

	return queryString
}

func replace(input, from string, to interface{}) string {
	s := fmt.Sprintf("%v", to)
	return strings.Replace(input, from, s, -1)
}

func fetchDataInRange(action, endpoint, request, token string) []interface{} {
	req, _ := b64.StdEncoding.DecodeString(request)
	var res []interface{}
	request = html.UnescapeString(string(req))
	// log.Println("fetchDataInRange", request)
	if action == "GET" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, json2query(request))
		// log.Println(endpoint)
		r := requestGetAPI(endpoint, request, token)
		res = r
	}

	if action == "POST" {
		r := requestPostAPI(endpoint, request, token)
		res = r
	}
	// log.Println("fetchDataInRange", res)
	return res
}

func requestGetAPI(endpoint, request, token string) []interface{} {
	var res APIResponse
	requestData := bytes.NewBuffer([]byte(request))

	req, err := http.NewRequest("GET", endpoint, requestData)
	if err != nil {
		log.Println("requestGetAPI Request err:", err)
		return nil
	}

	// 暫時還沒有加入jwt
	req.Header.Set("CGA-Header", "cga-good-good")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Bearer", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err)
		return nil
	}

	return res.Data
}

func requestPostAPI(endpoint, request, token string) []interface{} {
	var res APIResponse
	requestData := bytes.NewBuffer([]byte(request))

	req, err := http.NewRequest("POST", endpoint, requestData)
	if err != nil {
		log.Println("requestGetAPI Request err:", err)
		return nil
	}

	// 暫時還沒有加入jwt
	req.Header.Set("CGA-Header", "cga-good-good")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Bearer", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err)
		return nil
	}

	return res.Data
}

func getNewHtml(temples string) string {
	funcMap := template.FuncMap{
		"FetchData": fetchDataInRange,
		"Replace":   replace,
	}
	// log.Println("temples", temples)
	t, err := template.New("tmp").Funcs(funcMap).Parse(temples)
	if err != nil {
		log.Println("getNewHtml template err: ", err)
		return err.Error()
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, nil); err != nil {
		log.Println("getNewHtml template Execute err: ", err.Error())
		return err.Error()
	}

	// log.Println("temples", tpl.String())

	return tpl.String()
}

func ifTagToTmpl(content, logic string) string {
	if logic == "" {
		return content
	}
	tmplTop := fmt.Sprintf(`{{ if %s }}`, logic)
	tmplContent := content
	tmplEnd := `{{ end }}`
	tmpl := tmplTop + tmplContent + tmplEnd
	return tmpl
}

func forVarTagToTmpl(content, aliasName string) string {
	if aliasName == "" {
		return content
	}
	tmplTop := fmt.Sprintf(`{{ range $%sIndex, $%s := $%ss }}`, aliasName, aliasName, aliasName)
	tmplContent := content
	tmplEnd := `{{ end }}`
	tmpl := tmplTop + tmplContent + tmplEnd
	return tmpl
}

func rangeTagToTmpl(content, aliasName, endpoint, request, action, token string) string {
	if aliasName == "" {
		return content
	}

	endpointTmp := fmt.Sprintf(`{{ $%sendpoint := "%s" }}`, aliasName, endpoint)
	request = b64.StdEncoding.EncodeToString([]byte(request))
	requestTmp := fmt.Sprintf(`{{ $%srequest := "%s" }}`, aliasName, request)
	actionTmp := fmt.Sprintf(`{{ $%saction := "%s" }}`, aliasName, action)
	fetchData := fmt.Sprintf(`{{ $%ss := FetchData $%saction $%sendpoint $%srequest "%s" }}`, aliasName, aliasName, aliasName, aliasName, token)
	tmplVar := endpointTmp + requestTmp + actionTmp + fetchData
	tmplTop := fmt.Sprintf(`{{ range $index, $%s := $%ss }}`, aliasName, aliasName)
	tmplContent := content
	tmplEnd := `{{ end }}`
	tmpl := tmplVar + tmplTop + tmplContent + tmplEnd
	return tmpl
}

func fetchDataTagToTmpl(aliasName, endpoint, request, action, token string) string {
	if aliasName == "" || endpoint == "" {
		return ""
	}

	endpointTmp := fmt.Sprintf(`{{ $%sendpoint := "%s" }} {{ $%sendpoint }} `, aliasName, endpoint, aliasName)
	request = b64.StdEncoding.EncodeToString([]byte(request))
	requestTmp := fmt.Sprintf(`{{ $%srequest := "%s" }}`, aliasName, request)
	actionTmp := fmt.Sprintf(`{{ $%saction := "%s" }}`, aliasName, action)
	tmplVar := fmt.Sprintf(`{{ $%ss := FetchData $%saction $%sendpoint $%srequest "%s" }}`, aliasName, aliasName, aliasName, aliasName, token)
	tmplVarTotal := fmt.Sprintf(`{{ $%sTotal := len $%ss }}`, aliasName, aliasName)
	tmpl := endpointTmp + requestTmp + actionTmp + tmplVar + tmplVarTotal
	//log.Println(tmpl)
	return tmpl
}
