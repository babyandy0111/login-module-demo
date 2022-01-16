package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := "$test.id/$qq.test"
	regexp, _ := regexp.Compile(`\$[a-zA-Z0-9]*[a-zA-Z0-9]`)
	test := regexp.FindAllString(str, -1)
	fmt.Println(test)
}

//import (
//	"fmt"
//	"github.com/PuerkitoBio/goquery"
//	"log"
//	"os"
//)
//
//func main() {
//	f, err := os.Open("test.html")
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	defer f.Close()
//	dom, err := goquery.NewDocumentFromReader(f)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	dom.Find("range").Each(func(i int, s *goquery.Selection) {
//		aliasName := s.AttrOr("data-gjs-aliasname", "")
//		content, err := s.Html()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		s.ReplaceWithHtml(rangeTagToTmpl(content, aliasName))
//	})
//
//	dom.Find("for-var").Each(func(i int, s *goquery.Selection) {
//		aliasName := s.AttrOr("data-gjs-varname", "")
//		content, err := s.Html()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		s.ReplaceWithHtml(forVarTagToTmpl(content, aliasName))
//	})
//
//	dom.Find("if").Each(func(i int, s *goquery.Selection) {
//		logic := s.AttrOr("data-gjs-logic", "")
//		content, err := s.Html()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		s.ReplaceWithHtml(ifTagToTmpl(content, logic))
//	})
//
//	dom.Find("fetchdata").Each(func(i int, s *goquery.Selection) {
//		aliasName := s.AttrOr("data-gjs-aliasname", "")
//		endpoint := s.AttrOr("data-gjs-endpoint", "")
//		request := s.AttrOr("data-gjs-request", "{&quot;menu.pid&quot;:&quot;$menu.id&quot;}")
//		action := s.AttrOr("data-gjs-action", "GET")
//
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		s.ReplaceWithHtml(fetchDataTagToTmpl(aliasName, endpoint, request, action))
//	})
//
//	html, err := dom.Html()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	fmt.Printf("%+v", html)
//}
//
//func ifTagToTmpl(content, logic string) string {
//	if logic == "" {
//		return content
//	}
//
//	tmplTop := `{{ if %s }}`
//	tmplTop = fmt.Sprintf(tmplTop, logic)
//	tmplContent := content
//	tmplEnd := `{{ end }}`
//	tmpl := tmplTop + tmplContent + tmplEnd
//	return tmpl
//}
//
//func forVarTagToTmpl(content, aliasName string) string {
//	if aliasName == "" {
//		return content
//	}
//	tmplTop := `{{ range $%sIndex, $%s := $%ss }}`
//	tmplTop = fmt.Sprintf(tmplTop, aliasName, aliasName, aliasName)
//	tmplContent := content
//	tmplEnd := `{{ end }}`
//	tmpl := tmplTop + tmplContent + tmplEnd
//	return tmpl
//}
//
//func rangeTagToTmpl(content, aliasName string) string {
//	if aliasName == "" {
//		return content
//	}
//	tmplVar := `{{ $%ss := .%s }}`
//	tmplVar = fmt.Sprintf(tmplVar, aliasName, aliasName)
//	tmplTop := `{{ range $index, $%s := $%ss }}`
//	tmplTop = fmt.Sprintf(tmplTop, aliasName, aliasName)
//	tmplContent := content
//	tmplEnd := `{{ end }}`
//	tmpl := tmplVar + tmplTop + tmplContent + tmplEnd
//	return tmpl
//}
//
//func fetchDataTagToTmpl(aliasName, endpoint, request, action string) string {
//	if aliasName == "" || endpoint == "" {
//		return ""
//	}
//	endpointTmp := `{{ $%sendpoint := Replace "%s" "$%s" $%s }}`
//	endpointTmp = fmt.Sprintf(endpointTmp, aliasName, endpoint, aliasName, aliasName)
//
//	requestTmp := `{{ $%srequest := Replace "%s" "$%s" $%s }}`
//	requestTmp = fmt.Sprintf(requestTmp, aliasName, request, aliasName, aliasName)
//
//	actionTmp := `{{ $%saction := "%s" }}`
//	actionTmp = fmt.Sprintf(actionTmp, aliasName, action)
//
//	varP := `{{ $%ss := FetchData $%saction $%sendpoint $%srequest "" }}`
//	varP = fmt.Sprintf(varP, aliasName, aliasName, aliasName, aliasName)
//	varPTotal := `{{ $%sTotal := len $%ss }}`
//	varPTotal = fmt.Sprintf(varPTotal, aliasName, aliasName)
//
//	tmpl := endpointTmp + requestTmp + actionTmp + varP + varPTotal
//	return tmpl
//}
