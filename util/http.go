package util

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Paginate(header http.Header, count int, page int, perPage int, uri string) {
	totalPages := int(math.Ceil(float64(count / perPage)))
	lastPage := totalPages
	curPage := page
	prevPage := 0
	if curPage >= 2 {
		prevPage = curPage - 1
	}
	nextPage := 0
	if (curPage + 1) <= totalPages {
		nextPage = curPage + 1
	}

	header.Set("X-Page", strconv.Itoa(curPage))
	header.Set("X-Per-Page", strconv.Itoa(perPage))
	header.Set("X-Total", strconv.Itoa(count))
	header.Set("X-Total-Pages", strconv.Itoa(totalPages))

	if nextPage > 0 {
		header.Set("X-Next-Page", strconv.Itoa(nextPage))
		appendLink(header, formatURI(uri, perPage, nextPage), "next", nil)
	}

	appendLink(header, formatURI(uri, perPage, lastPage), "last", nil)
	appendLink(header, formatURI(uri, perPage, 1), "first", nil)

	if prevPage > 0 {
		header.Set("X-Prev-Page", strconv.Itoa(prevPage))
		appendLink(header, formatURI(uri, perPage, prevPage), "prev", nil)
	}
}

func formatURI(uri string, perPage int, page int) string {
	return fmt.Sprintf("%s?per_page=%d&page=%d", uri, perPage, page)
}

type LinkOptions struct {
	Title         string
	TitleStar     string
	Anchor        string
	HREFLang      []string
	TypeHint      string
	CrossOrigin   string
	LinkExtension [][]string
}

func appendLink(header http.Header, target string, rel string, opts *LinkOptions) {
	if strings.Contains(rel, "//") {
		if strings.Contains(rel, " ") {
			var escaped []string
			for _, s := range strings.Split(rel, " ") {
				escaped = append(escaped, url.QueryEscape(s))
			}
			rel = fmt.Sprintf(`"%s"`, strings.Join(escaped, " "))
		} else {
			rel = fmt.Sprintf(`"%s"`, url.QueryEscape(rel))
		}
	}

	// value := fmt.Sprintf("<%s>; rel=%s", url.QueryEscape(target), rel)
	value := fmt.Sprintf("<%s>; rel=%s", target, rel)

	if opts != nil {
		if opts.Title != "" {
			value += fmt.Sprintf(`; title="%s"`, opts.Title)
		}

		if opts.TitleStar != "" {
			value += fmt.Sprintf(`; title=*UTF-8'%s'%s`, opts.TitleStar[0], url.QueryEscape(opts.TitleStar[:1]))
		}

		if opts.TypeHint != "" {
			value += fmt.Sprintf(`; type="%s"`, opts.TypeHint)
		}

		if len(opts.HREFLang) > 0 {
			var langs []string
			for _, s := range opts.HREFLang {
				langs = append(langs, "hreflang="+s)
			}
			value += "; "
			value += strings.Join(langs, "; ")
		}

		if opts.Anchor != "" {
			value += fmt.Sprintf(`; achor="%s"`, url.QueryEscape(opts.Anchor))
		}

		if opts.CrossOrigin != "" {
			opts.CrossOrigin = strings.ToLower(opts.CrossOrigin)
			if opts.CrossOrigin == "anonymous" {
				value += "; crossorigin"
			} else {
				value += `; crossorigin="use-credentials"`
			}
		}

		if len(opts.LinkExtension) > 0 {
			var links []string
			for _, s := range opts.LinkExtension {
				links = append(links, fmt.Sprintf("%s=%s", s[0], s[1]))
			}
			value += "; "
			value += strings.Join(links, "; ")
		}
	}

	link := header.Get("Link")
	if link != "" {
		link += fmt.Sprintf(", %s", value)
		header.Set("Link", link)
	} else {
		header.Set("Link", value)
	}
}
