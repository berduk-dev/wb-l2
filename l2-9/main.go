package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("неверный запрос")
		return
	}
	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error http.Get(url):", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	file, err := os.Create("index.html")
	if err != nil {
		fmt.Println("error os.Create:", err)
	}

	_, err = file.Write(body)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

}

func ExtractLinks(baseURL string, r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var links []string
	var visit func(*html.Node)

	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// Преобразуем относительные ссылки в абсолютные
					u, err := url.Parse(attr.Val)
					if err == nil {
						base, _ := url.Parse(baseURL)
						link := base.ResolveReference(u)
						links = append(links, link.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}

	visit(doc)
	return links, nil
}

var visited = make(map[string]bool)

func Download(url string) {
	// Если мы уже были на этой странице — выходим
	if visited[url] {
		return
	}
	visited[url] = true

	fmt.Println("Скачиваю:", url)

	// Делаем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	// Извлекаем все ссылки со страницы
	links, err := ExtractLinks(url, resp.Body)
	if err != nil {
		fmt.Println("Ошибка парсинга:", err)
		return
	}

	// Рекурсивно обходим найденные ссылки
	for _, link := range links {
		Download(link)
	}
}
