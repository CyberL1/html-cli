package utils

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func ApplyBoilerplate(fileBytes []byte, includeHotReload bool) []byte {
	contents := string(fileBytes)
	document, err := html.Parse(strings.NewReader(contents))
	if err != nil {
		fmt.Println("Failed to parse HTML:", err)
		return nil
	}

	head := findElement(document, "head")

	if includeHotReload {
		head.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "script",
			Attr: []html.Attribute{
				{Key: "id", Val: "html-cli-hot-reload"},
			},
		})

		findElementWithAttr(head, "script", "id", "html-cli-hot-reload").AppendChild(&html.Node{
			Type: html.TextNode,
			Data: "(()=>{const sse=new EventSource('/_html/hot-reload');sse.onopen=()=>{console.log('* HTML hot-reload enabled *')};sse.onmessage=e=>{if(e.data==='reload'){console.log('* HTML hot-reload triggered *');location.reload()}}})()",
		})
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, document); err != nil {
		fmt.Println("Failed to render HTML:", err)
		return nil
	}
	return buf.Bytes()
}
