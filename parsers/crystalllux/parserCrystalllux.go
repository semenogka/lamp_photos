package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type SiteStruct struct {
	URL      string
	Category string
	Pages    int
}

type Product struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	Height    string `json:"height"`
	Width    string `json:"width"`
	Length    string `json:"length"`
	Diameter    string `json:"diameter"`
	ArmaturColor    string `json:"metalcolor"`
	AbajurColor string `json:"abajurcolor"`
	ArmaturMaterial    string `json:"armaturmaterial"`
	AbajurMaterial string `json:"abajurmaterial"`
}

var crystalluxSites = []SiteStruct{
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/lyustry/?PAGEN_1=", Category: "alldata\\data\\lystri.json", Pages: 26},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/torshery/?PAGEN_1=", Category: "alldata\\data\\torsher.json", Pages: 2},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/bra/?PAGEN_1=", Category: "alldata\\data\\bra.json", Pages: 22},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/nastolnye-lampy/?PAGEN_1=", Category: "alldata\\data\\nastol.json", Pages: 3},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/svetilniki/?PAGEN_1=", Category: "alldata\\data\\podves.json", Pages: 24},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/dekorativnyy-svet-potolochnye-svetilniki/?PAGEN_1=", Category: "alldata\\data\\potol.json", Pages: 5},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/clt-nastennye-svetilniki-bra/?PAGEN_1=", Category: "alldata\\data\\bra.json", Pages: 12},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/vstraivaemye/?PAGEN_1=", Category: "alldata\\data\\tochvstr.json", Pages: 13},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/trekovye-sistemy/svetilniki2/?PAGEN_1=", Category: "alldata\\data\\magn.json", Pages: 6},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/clt-potolochnye-svetilniki/?PAGEN_1=", Category: "alldata\\data\\potol.json", Pages: 6},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/odnofaznaya-trekovaya-sistema/odnofaznaya-svetilniki/?PAGEN_1=", Category: "alldata\\data\\track.json", Pages: 4},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/technical-clt-podvesnye-svetilniki/?PAGEN_1=", Category: "alldata\\data\\tochpodv.json", Pages: 6},
}


func main() {
	//links := []string{}
	for _, site := range crystalluxSites {
		var products []Product
		links, names := takeLinks(site.URL, site.Pages)
		
		for i, link := range links {
			c := colly.NewCollector()
			product := Product{
				Height:          "None",
				Width:          "None",
				Length: 		 "None",
				Diameter:        "None",
				ArmaturColor:    "None",
				AbajurColor:     "None",
				ArmaturMaterial: "None",
				AbajurMaterial:  "None",
			}
			
			c.OnHTML(".new-product-item-detail-properties", func(e *colly.HTMLElement) {
				product.Link = link
				product.Name = names[i]
				var maxHeight int
				var length, width string
				e.DOM.Find("dt").Each(func(i int, s *goquery.Selection) {
				text := s.Text()
				text = strings.TrimSpace(text)
				
				next := s.Next()
				nextText := next.Text()
				
				if strings.HasPrefix(text, "Высота изделия, мм") {
					numText := maxValue(nextText)
					
					if  numText > maxHeight {
						maxHeight = numText
						product.Height = strconv.Itoa(maxHeight)
					}
				} else if strings.HasPrefix(text, "Длина") && !strings.HasSuffix(text, "коробки"){
					if length == "" { // берем только первый
						num := maxValue(nextText)
						length = strconv.Itoa(num)
						product.Length = length
					}
				} else if strings.HasPrefix(text, "Ширина") {
					if width == "" { // берем только первый
						num := maxValue(nextText)
						width = strconv.Itoa(num)
						product.Width = width
					}
				}else if strings.HasPrefix(text, "Диаметр") {
					if product.Diameter == "None" { // берем только первый
						num := maxValue(nextText)
						product.Diameter = strconv.Itoa(num)
					}
				} else if strings.HasPrefix(text, "Материал арматуры") {
					value := cleanText(nextText)
					product.ArmaturMaterial = replaceSeparatorsWithSlash(value)
					
				} else if strings.HasPrefix(text, "Цвет арматуры") {
					value := cleanText(nextText)
					product.ArmaturColor = replaceSeparatorsWithSlash(value)
				} else if strings.HasPrefix(text, "Цвет абажура") {
					value := cleanText(nextText)
					product.AbajurColor = replaceSeparatorsWithSlash(value)
				} else if strings.HasPrefix(text, "Материал абажура") {
					value := cleanText(nextText)
					product.AbajurMaterial = replaceSeparatorsWithSlash(value)
				}
			})
		})

		c.OnScraped(func(r *colly.Response) {
			products = append(products, product)
			log.Println(product)
		}) 

		c.Visit(link)
			

		
		}
		saveToJson(site.Category, products)
	}
	
}

func downloadImg(src string, name string) {
	filepath := filepath.Join("allimgs", name)

	if _, err := os.Stat(filepath); err == nil {
		// Файл уже существует — не скачиваем
		log.Println("Файл уже существует, пропускаем:", name)
		return
	}

	res, err := http.Get(src)

	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	

	out, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Println(err)
	}
	
}

func saveToJson(filename string, data interface{}) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	var exData []interface{}

	stat, _ := file.Stat()
	if stat.Size() != 0 {
		json.NewDecoder(file).Decode(&exData)
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			exData = append(exData, v.Index(i).Interface())
		}
	} else {
		exData = append(exData, data)
	}

	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Ошибка seek:", err)
		return
	}
	if err := file.Truncate(0); err != nil {
		log.Println("Ошибка truncate:", err)
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	err = encoder.Encode(exData)
	if err != nil {
		log.Println(err)
	}
}

func takeLinks(link string, pages int) ([]string, []string) {
	var links []string
	var names []string
	
	c := colly.NewCollector()

	
		c.OnHTML(".product-item", func(e *colly.HTMLElement) {
			link := e.ChildAttr("a", "href")
			fullLink := e.Request.AbsoluteURL(link)
			//log.Println("Ссылка на товар:", fullLink)

			
			imgUrl := e.ChildAttr("a meta", "content")
			fullImgURL := e.Request.AbsoluteURL(imgUrl)
			filename := fmt.Sprintf("CRYSTALLUX%s.jpg", extractID(fullImgURL))
			
			downloadImg(fullImgURL, filename)
			links = append(links, fullLink)
			names = append(names, filename)
			

			

			//log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename)
		})	

		for i := 1; i <= pages; i++{
			url := "%s%d"


			page := fmt.Sprintf(url, link, i)
			log.Println(page)
			err := c.Visit(page)
			if err != nil {
				log.Fatal(err)
			}
		}

	return links, names
}

func extractNumber(s string) string {
	re := regexp.MustCompile(`\d+`)
	return strings.Join(re.FindAllString(s, -1), "")
}



func cleanValue(value string) string {
    value = strings.TrimSpace(value)
    parts := strings.Split(value, ",")
    return strings.TrimSpace(parts[0])
}

func cleanText(raw string) string {
    // Заменяем неразрывный пробел на обычный пробел
    cleaned := strings.ReplaceAll(raw, "\u00A0", " ")
    // Убираем дефис и лишние пробелы
    cleaned = strings.ReplaceAll(cleaned, "—", "")
    // Удаляем все лишние пробелы
    cleaned = strings.TrimSpace(cleaned)
    return cleaned
}

func extractImgUrl(style string) string{
	runes := []rune{}
	count := false
	var s string
	for _, r := range style {
		if r == '(' {
			count = true
		}
		if count && r != '('{
			runes = append(runes, r)
		}
		if len(runes) >= 4 {
			s = string(runes[:len(runes)-3]) // удаляем последние 3
			s = s[1:]    // удаляем первый
		}
	}
	return s
}


func extractName(url string) string {
	count := 0
	name := []rune{}
	for _, let := range url {
		if let == '/' {
			count += 1
		}

		if count == 6 && let != '/'{
			name = append(name, let)
		}
	}
	filename := fmt.Sprintf("%s.jpg", string(name))
	return filename
}

func extractID(url string) (string) {
    re := regexp.MustCompile(`/upload/iblock/[0-9a-z]+/([0-9a-z]+)`)
    match := re.FindStringSubmatch(url)
    if len(match) > 1 {
        return match[1]
    }
    return ""
}

func maxValue(s string) int {
	parts := strings.Split(s, "/")
	if len(parts) == 1 {
		parts = strings.Split(s, "-")
		if len(parts) == 1 {
			parts = strings.Split(s, "–")
			if len(parts) == 1 {
				parts = strings.Split(s, "—")
			}
		}
	}
	max := 0
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.Atoi(extractNumber(part))
		
		if err == nil && num > max {
			max = num
		}
	}
	return max
}

func replaceSeparatorsWithSlash(s string) string {
	reSep := regexp.MustCompile(`\s*(,|;)\s*`)
	s = reSep.ReplaceAllString(s, "/")

	reSlash := regexp.MustCompile(`\s*/\s*`)
	return reSlash.ReplaceAllString(s, "/")
}