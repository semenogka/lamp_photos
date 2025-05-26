package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

type Product struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
}

var products []Product

var pages int = 15



func main() {
	c := colly.NewCollector()
	
	c.OnHTML(".product-card-inner", func(h *colly.HTMLElement) {
		title := h.ChildText(".flex-1 .px-2 .product-card-title div")
		if strings.Contains(title, "Светильник"){
			url := h.ChildAttr(".flex-1 .d-flex .top-icons-box .copy-link", "data-url")
			imgUrl := h.ChildAttr(".flex-1 .product-card-image .position-relative .list-card-slider img", "src")
			fullImgURL := h.Request.AbsoluteURL(imgUrl)
			filename := fmt.Sprintf("ISONEX%s", extractName(fullImgURL))

			products = append(products, Product{
				Name: filename,
				Link: url,
			})

			downloadImg(fullImgURL, filename)
			log.Println(url, " ", fullImgURL, " ", filename, " ", len(products))
		}
		
	})
	for i := 1; i <= pages; i++{
		page := fmt.Sprintf("https://isonex.ru/catalog/filter/dev_site_category-is-%%D0%%B2%%D1%%81%%D1%%82%%D1%%80%%D0%%B0%%D0%%B8%%D0%%B2%%D0%%B0%%D0%%B5%%D0%%BC%%D1%%8B%%D0%%B5%%20%%D1%%81%%D0%%B2%%D0%%B5%%D1%%82%%D0%%B8%%D0%%BB%%D1%%8C%%D0%%BD%%D0%%B8%%D0%%BA%%D0%%B8/apply/?ajax_request=Y&PAGEN_1=%d", i)

		log.Println(page)
		c.Visit(page)
	}
	
	
	saveToJson("alldata\\data\\tochvstr.json", products)
}

func downloadImg(src string, name string) {
	res, err := http.Get(src)

	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	filepath := filepath.Join("allimgs", name)

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

	exData = append(exData, data)


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

func extractName(url string) string {
	count := 0
	name := []rune{}
	for _, let := range url {
		if let == '/' {
			count += 1
		}

		if count == 8 && let != '/'{
			name = append(name, let)
		}
	}
	
	return string(name)
}
