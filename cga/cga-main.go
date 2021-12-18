package cga

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"net/http"
)

var engine *xorm.Engine

type Blocks struct {
	Id         string           `json:"id"`
	Media      string           `json:"media"`
	Label      string           `json:"label"`
	Category   string           `json:"category"`
	Attributes BlocksAttributes `xorm:"json" json:"attributes"`
	Content    string           `json:"content"`
}

type BlocksAttributes struct {
	Class string `json:"class"`
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.LoadHTMLGlob("templates/view/*.html")
	router.Static("/statics", "templates/statics")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
			"Langs": []string{"Python", "Ruby", "PHP", "Java", "Golang"},
		})
	})

	router.GET("/test", func(c *gin.Context) {
		var jsonStr = `[{"id":"section-brands-test","media": "","label":"Brands-test","category": "Sections","attributes": {"class": "fa fa-ellipsis-h"},"content": "<section id=\"brands\" class=\"bg-white\"><div class=\"container\"><div class=\"row\">{{ range .TestData }}<div class=\"col-lg-2 col-md-3 text-center\"><p class=\"text-muted mb-2\">{{ . }}</p><img src=\"https:\/\/upload.wikimedia.org\/wikipedia\/commons\/5\/53\/Google_%22G%22_Logo.svg\" alt=\"Google\"></div>{{ end }}</div></div></section><style>#brands img{opacity:0.3;}</style>"},{"id":"section-brands-test2","media": "","label":"Brands-test2","category": "Sections","attributes": {"class": "fa fa-ellipsis-h"},"content": "<section id=\"brands\" class=\"bg-white\"><div class=\"container\"><div class=\"row\">{{ range .TestData }}<div class=\"col-lg-2 col-md-3 text-center\"><p class=\"text-muted mb-2\">{{ . }}</p><img src=\"https:\/\/upload.wikimedia.org\/wikipedia\/commons\/5\/53\/Google_%22G%22_Logo.svg\" alt=\"Google\"></div>{{ end }}</div></div></section><style>#brands img{opacity:0.3;}</style>"}]`
		jsonParsed, err := gabs.ParseJSON([]byte(jsonStr))
		if err != nil {
			panic(err)
		}

		c.JSON(200, jsonParsed)
	})

	router.GET("/blocks", func(c *gin.Context) {
		result, _ := blockData()

		jsonParsed, err := gabs.ParseJSON([]byte(result))
		if err != nil {
			panic(err)
		}

		c.JSON(200, jsonParsed)
	})

	router.Run(":8080")
}

func blockData() ([]byte, error) {
	//连接数据库
	engine, err := xorm.NewEngine("mysql", "root:8eps5tEtWuwr@tcp(preview-project.cpzq1quzpyg0.us-east-1.rds.amazonaws.com:3306)/automation?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer engine.Close() //延迟关闭数据库
	// fmt.Println("数据库链接成功")

	//查询单条数据
	// blocksList := make([]Blocks, 0)
	//sql := "SELECT json_object('id', id, 'lable', lable, 'category', category, 'attributes', attributes, 'content', content) as result FROM `blocks`"
	sql := "SELECT blocks_id as id, lable, category, content, attributes FROM `blocks`"
	// results, err := engine.QueryString("SELECT json_object('id', id, 'lable', lable, 'category', category, 'attributes', attributes, 'content', content) as result FROM `blocks`")

	result, err := engine.QueryResult(sql).List()

	// blocks := make([string]int, len(result))
	var blocks []interface{}
	for _, data := range result {
		var attributes BlocksAttributes
		json.Unmarshal(data["attributes"], &attributes)

		resp := Blocks{
			Id:       string(data["id"]),
			Media:    string(data["media"]),
			Label:    string(data["lable"]),
			Content:  string(data["content"]),
			Category: string(data["category"]),
			Attributes: BlocksAttributes{
				Class: attributes.Class,
			},
		}

		blocks = append(blocks, resp)
	}
	b, err := json.Marshal(blocks)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b, err
}
