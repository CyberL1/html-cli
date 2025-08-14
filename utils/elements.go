package utils

import "golang.org/x/net/html"

func findElement(document *html.Node, element string) *html.Node {
	if document.Type == html.ElementNode && document.Data == element {
		return document
	}

	for c := document.FirstChild; c != nil; c = c.NextSibling {
		if result := findElement(c, element); result != nil {
			return result
		}
	}
	return nil
}

func findElementWithAttr(document *html.Node, element, attrKey, attrVal string) *html.Node {
	if document.Type == html.ElementNode && document.Data == element {
		for _, attr := range document.Attr {
			if attr.Key == attrKey && attr.Val == attrVal {
				return document
			}
		}
	}

	for c := document.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementWithAttr(c, element, attrKey, attrVal); result != nil {
			return result
		}
	}
	return nil
}
