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
	"regexp"
	"strings"
)

type APIResponse struct {
	Data  []interface{} `json:"data"`
	Total int           `json:"total"`
}

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {

		token := ""
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
		request := s.AttrOr("data-gjs-request", "{&#34;xxx&#34;:&#34;0&#34;}")
		action := s.AttrOr("data-gjs-action", "GET")

		content, err := s.Html()
		if err != nil {
			log.Println(err)
			return
		}
		s.ReplaceWithHtml(rangeTagToTmpl(content, aliasName, endpoint, request, action, token))
	})

	dom.Find("for").Each(func(i int, s *goquery.Selection) {
		aliasName := s.AttrOr("data-gjs-varname", "")
		content, err := s.Html()
		if err != nil {
			log.Println(err)
			return
		}
		s.ReplaceWithHtml(forTagToTmpl(content, aliasName))
	})

	dom.Find("if").Each(func(i int, s *goquery.Selection) {
		logic := s.AttrOr("data-gjs-logic", "gt 1 0")
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
		request := s.AttrOr("data-gjs-request", "{&#34;xxx&#34;:&#34;0&#34;}")
		action := s.AttrOr("data-gjs-action", "GET")
		s.ReplaceWithHtml(fetchDataTagToTmpl(aliasName, endpoint, request, action, token))
	})

	dom.Find("member").Each(func(i int, s *goquery.Selection) {
		content, err := s.Html()
		if err != nil {
			log.Println(err)
			return
		}
		s.ReplaceWithHtml(memberTagToTmpl(content, token))
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

func replaceJson(values ...interface{}) string {
	if len(values)%2 != 0 {
		log.Println("len(values)%2 err")
		return values[0].(string)
	}

	source, err := b64.StdEncoding.DecodeString(values[0].(string))
	if err != nil {
		return values[0].(string)
	}
	// log.Println("source", source)
	jsonRes := string(source)
	//dict := make(map[string]interface{}, len(values)/2)
	for i := 2; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return "dict keys must be strings"
		}
		s := fmt.Sprintf("%v", values[i+1])
		jsonRes = strings.Replace(jsonRes, key, s, -1)
		//dict[key] = values[i+1]
		// log.Println("loop:", key, values[i+1])
	}
	// log.Println("replaceJson:", jsonRes)
	return b64.StdEncoding.EncodeToString([]byte(jsonRes))
}

func replaceEndpoint(values ...interface{}) string {
	if len(values)%2 != 0 {
		log.Println("len(values)%2 err")
		return values[0].(string)
	}

	source := values[0].(string)

	//dict := make(map[string]interface{}, len(values)/2)
	for i := 2; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return "dict keys must be strings"
		}
		s := fmt.Sprintf("%v", values[i+1])
		source = strings.Replace(source, key, s, -1)
		//dict[key] = values[i+1]
		// log.Println("loop:", key, values[i+1])
	}
	// log.Println("source", source)
	return html.UnescapeString(source)
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
		"FetchData":       fetchDataInRange,
		"ReplaceEndpoint": replaceEndpoint,
		"ReplaceJson":     replaceJson,
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

func forTagToTmpl(content, aliasName string) string {
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
	//endpointTmp := fmt.Sprintf(`{{ $%sendpoint := "%s" }}`, aliasName, endpoint)
	//request = b64.StdEncoding.EncodeToString([]byte(request))
	//requestTmp := fmt.Sprintf(`{{ $%srequest := "%s" }}`, aliasName, request)
	//actionTmp := fmt.Sprintf(`{{ $%saction := "%s" }}`, aliasName, action)
	//fetchData := fmt.Sprintf(`{{ $%ss := FetchData $%saction $%sendpoint $%srequest "%s" }}`, aliasName, aliasName, aliasName, aliasName, token)
	// tmplVar := endpointTmp + requestTmp + actionTmp + fetchData
	fetchData := fetchDataTagToTmpl(aliasName, endpoint, request, action, token)
	tmplTop := fmt.Sprintf(`{{ range $index, $%s := $%ss }}`, aliasName, aliasName)
	tmplContent := content
	tmplEnd := `{{ end }}`
	tmpl := fetchData + tmplTop + tmplContent + tmplEnd
	return tmpl
}

func fetchDataTagToTmpl(aliasName, endpoint, request, action, token string) string {
	if aliasName == "" || endpoint == "" {
		return ""
	}

	regexp, _ := regexp.Compile(`\$[a-zA-Z0-9]*[a-zA-Z0-9]\.[a-zA-Z0-9]*[a-zA-Z0-9]`)
	resEndpoint := regexp.FindAllString(endpoint, -1)
	tmpEndpointVar := ""
	for _, v := range resEndpoint {
		tmpEndpointVar = tmpEndpointVar + "\"" + v + "\"" + " " + v + " "
	}

	resJson := regexp.FindAllString(request, -1)
	tmpRequestVar := ""
	for _, v := range resJson {
		tmpRequestVar = tmpRequestVar + "\"" + v + "\"" + " " + v + " "
	}

	request = b64.StdEncoding.EncodeToString([]byte(request))

	// log.Println("tmpVar", tmpRequestVar)

	endpointTmp := fmt.Sprintf(`{{ $%sendpoint := ReplaceEndpoint "%s" "" %s }}`, aliasName, endpoint, tmpEndpointVar)

	requestTmp := fmt.Sprintf(`{{ $%srequest := ReplaceJson "%s" "" %s }}`, aliasName, request, tmpRequestVar)

	actionTmp := fmt.Sprintf(`{{ $%saction := "%s" }}`, aliasName, action)
	tmplVar := fmt.Sprintf(`{{ $%ss := FetchData $%saction $%sendpoint $%srequest "%s" }}`, aliasName, aliasName, aliasName, aliasName, token)
	tmplVarTotal := fmt.Sprintf(`{{ $%sTotal := len $%ss }}`, aliasName, aliasName)
	tmpl := endpointTmp + requestTmp + actionTmp + tmplVar + tmplVarTotal
	//log.Println(tmpl)
	return tmpl
}

func memberTagToTmpl(content, token string) string {
	if token == "" {
		return ""
	}
	// todo check token
	if true {
		return content
	} else {
		return ""
	}
}
