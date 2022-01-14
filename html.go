package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type APIResponse struct {
	Data  []interface{} `json:"data"`
	Total int           `json:"total"`
}

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		info := make(map[string]interface{})
		token := ""
		info["menu1"] = requestGetAPI("https://lmd4le8g.codegenapps.com/menu", "{&quot;menu.pid&quot;:&quot;0&quot;}", token)
		// info["menu2"] = getData("https://lmd4le8g.codegenapps.com/menu?menu.pid=1")

		_ =
			`
			<html><head><meta charset="utf-8"></head>
			<ul class="navigation-menu">
			<range id="0705cf66-d77b-428d-bda3-ddb2d0db16e7" 
			data-gjs-aliasname="menu1" 
			data-gjs-endpoint="https://lmd4le8g.codegenapps.com/menu/$menu1.id" 
			data-gjs-request="{&quot;menu.pid&quot;:&quot;$menu.id&quot;}" 
			data-gjs-action="GET">
			  
				<li class="has-submenu parent-menu-item">
				  <a href="javascript:void(0)" id="ipcrp">找廠商</a>
				  
		
				  <featchdatainrange data-gjs-aliasname="menu2" 
					data-gjs-endpoint="https://lmd4le8g.codegenapps.com/menu/$menu1.id" 
					data-gjs-request="{&quot;menu.pid&quot;:&quot;$menu.id&quot;}" 
					data-gjs-action="GET">
		
		
				  <if id="22726769-fc0c-4cde-a1a8-25a304f7f9e4" gjs-data-logic="if gt xxxx 0">
					<span class="menu-arrow"></span>
					
					<rangeinrange id="0705cf66-d77b-428d-bda3-ddb2d0db16e7-2">
		
					  <ul class="submenu">
						<li>
						  <a href="" class="sub-menu-item">
							平面設計公司
						  </a>
						</li>
					  </ul>
					</rangeinrange>
		
				  </if>
		
		
				</li>
			 
			</range>
		 </ul>
			`

		tmpHtml :=
			`
			<html><head><meta charset="utf-8"></head>

			{{ $menu1s := .menu1 }}
			{{ range $index, $menu1 := $menu1s }}
				
				<li class="has-submenu parent-parent-menu-item">
					<a href="https://lmd4le8g.codegenapps.com/menu?menu.pid={{ $menu1.id }}" target="_blank">
						{{ $index }} - {{ $menu1.id }} - {{ $menu1.menu_name }}
					</a>
				</li>
				<br>

				{{ $endpoint := Replace "https://lmd4le8g.codegenapps.com/menu/$menu1.id" "menu1.id" $menu1.id }}
				{{ $request := Replace "{&quot;menu.pid&quot;:&quot;$menu1.id&quot;}" "$menu1.id" $menu1.id }}
				{{ $action := "GET" }}
				{{ $menu2s := FetchDataInRange $action $endpoint $request "" }}
				{{ $menu2Total := len $menu2s }}
				
				{{ if gt $menu2Total 0 }}
					有資料
					{{ range $index, $menu2 := $menu2s }}
						<a href="javascript:void(0)">
							{{ $menu2 }}
						</a>
						<br>
					{{ end }}
				{{ end }}

			{{ end }}
			</html>
			`

		// log.Println(editorHtml)
		html := getNewHtml(info, tmpHtml)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(html)
	})
	app.Listen(":3000")
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
	var res []interface{}
	request = html.UnescapeString(request)
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

func getNewHtml(apiInfo map[string]interface{}, temples string) string {
	funcMap := template.FuncMap{
		"FetchDataInRange": fetchDataInRange,
		"Replace":          replace,
	}
	t, err := template.New("tmp").Funcs(funcMap).Parse(temples)
	if err != nil {
		log.Println("getNewHtml template err: ", err)
		return err.Error()
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, apiInfo); err != nil {
		log.Println("getNewHtml template Execute err: ", err.Error())
		return err.Error()
	}

	return tpl.String()
}
