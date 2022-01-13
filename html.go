package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
		info["menu1"] = getData("https://lmd4le8g.codegenapps.com/menu?menu.pid=0")
		// info["menu2"] = getData("https://lmd4le8g.codegenapps.com/menu?menu.pid=1")

		// log.Println(info)

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

				{{ $endpoint := Replace "https://lmd4le8g.codegenapps.com/menu" "$menu1.id" $menu1.id }}
				{{ $request := Replace "{\"menu.pid\":\"$menu1.id\"}" "$menu1.id" $menu1.id }}
				{{ $action := "GET" }}
				{{ $menu2s := FetchData $menu1 $action $endpoint $request }}
				{{ range $index, $menu2 := $menu2s }}
					<a href="javascript:void(0)">
						{{ $menu2 }}
					</a>
					<br>
				{{ end }}
				
			{{ end }}
			</html>
			`

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

func gogo(mapInterface map[string]interface{}, action, endpoint, request string) []interface{} {
	if action == "GET" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, json2query(request))
	}
	// s := mapString["id"]
	// log.Println(mapInterface["id"])
	log.Println(endpoint, request)
	// url := fmt.Sprintf("https://lmd4le8g.codegenapps.com/menu?menu.pid=%v", mapInterface[0]["id"])
	// log.Println("https://lmd4le8g.codegenapps.com/menu?menu.pid=")
	// return getData(url)
	return nil
}
func getData(url string) []interface{} {
	// log.Println(url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var resData APIResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		log.Fatal(err)
	}
	return resData.Data
}

func getNewHtml(apiInfo map[string]interface{}, temples string) string {
	funcMap := template.FuncMap{
		"FetchData": gogo,
		"Replace":   replace,
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
