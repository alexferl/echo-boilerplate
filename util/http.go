package util

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/alexferl/httplink"
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
		httplink.Append(header, formatURI(uri, perPage, nextPage), "next")
	}

	httplink.Append(header, formatURI(uri, perPage, lastPage), "last")
	httplink.Append(header, formatURI(uri, perPage, 1), "first")

	if prevPage > 0 {
		header.Set("X-Prev-Page", strconv.Itoa(prevPage))
		httplink.Append(header, formatURI(uri, perPage, prevPage), "prev")
	}
}

func formatURI(uri string, perPage int, page int) string {
	return fmt.Sprintf("%s?per_page=%d&page=%d", uri, perPage, page)
}
