package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type HrefSlice []string

type Slice map[string]struct{}

func (ms Slice) ToSliceSafe(baseUrl string) HrefSlice {
	hrefSlice := make(HrefSlice, 0, len(ms))
	for key := range ms {
		hrefSlice = hrefSlice.safeAppend(key, baseUrl)
	}

	return hrefSlice
}

func (ms Slice) ToSlice() HrefSlice {
	hrefSlice := make(HrefSlice, 0, len(ms))
	for key := range ms {
		hrefSlice = append(hrefSlice, key)
	}

	return hrefSlice
}

func (s HrefSlice) safeAppend(value, baseUrl string) HrefSlice {
	if !strings.HasPrefix(value, strings.TrimSuffix(baseUrl, "/")) {
		return s
	}

	return append(s, value)
}

func Parse(r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)
	links := make(Slice)

	for _, node := range nodes {
		builtLink := buildLink(node)
		if builtLink == "" {
			continue
		}
		links[builtLink] = struct{}{}
	}

	return links.ToSlice(), nil
}

func buildLink(n *html.Node) string {
	var ret string

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret = attr.Val
			break
		}
	}
	return ret
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "base") {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func GetHrefs(r io.Reader, baseUrl string) (HrefSlice, error) {
	links, err := Parse(r)
	if err != nil {
		return nil, err
	}

	baseUrl = strings.TrimSuffix(baseUrl, "/")

	hrefMap := make(Slice)

	for _, l := range links {
		l = strings.TrimSuffix(l, "/")

		if isNotRelatedToBaseUrl(l, baseUrl) {
			continue
		}

		switch {
		case hasMoreThanBaseUrl(l, baseUrl) || (strings.HasPrefix(l, "http") && hasBaseUrl(l, baseUrl)):
			hrefMap[l] = struct{}{}
		case strings.HasSuffix(l, ".html"):
			hrefMap[baseUrl+"/"+strings.TrimPrefix(l, "/")] = struct{}{}
		case strings.HasPrefix(l, "/"):
			hrefMap[baseUrl+l] = struct{}{}
		}
	}

	return hrefMap.ToSliceSafe(baseUrl), nil
}

func hasMoreThanBaseUrl(link, baseUrl string) bool {
	if hasBaseUrl(link, baseUrl) {
		lenLink := len([]rune(link))
		lenBase := len([]rune(baseUrl))
		if lenLink > lenBase {
			return true
		}
	}
	return false
}

func hasBaseUrl(link, baseurl string) bool {
	return strings.HasPrefix(link, baseurl)
}

func isNotRelatedToBaseUrl(link, baseUrl string) bool {
	return strings.Contains(link, "#") || (strings.Contains(link, "https://") && !hasBaseUrl(link, baseUrl))
}
