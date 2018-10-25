package moni

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	NotCrawled = iota
	CrawlReady
	CrawlRequestSent
	CrawlResponseRecieved
	CrawlComplete
	CrawlErrored
	CrawlNotAllowed
)

// ===================================================================
type Page struct {
	URL string

	Content []byte
	Links   map[string]*Page
	Ignored map[string]int

	CrawlState int

	StatusCode int
	Err        error

	LastCrawled time.Time
	Start       time.Time
	Finish      time.Time
}

var (
	pages Pagemap
)

// Pagemap
// ********************************************************************
type Pagemap map[string]*Page

// String will represent the Page
// ====================================================================
func (p *Page) String() string {
	str := fmt.Sprintf("%s: lastcrawled: %s,  duration: %v links: %d ignored: %d\n", p.URL, p.LastCrawled, p.Finish, len(p.Links), len(p.Ignored))
	return str
}

func GetPages() Pagemap {
	if pages == nil {
		pages = make(Pagemap)

		st := GetStorage()
		if _, err := st.FetchObject("pages", &pages); err != nil {
			log.Debugf("Empty pages %v, creating ...", err)
			// TODO ~ make sure the error is NOT found

			pages = make(Pagemap)
			_, err := st.StoreObject("pages", pages)
			if err != nil {
				log.Errorf("Store: failed to create pages: %v ", err)
				return pages
			}
		}
	}
	return pages
}

func savePagemap() error {
	st := GetStorage()
	if _, err := st.StoreObject("pages", pages); err != nil {
		log.Errorf("failed to save page map %v", err)
		return err
	}
	return nil
}

// GetPage will sanitize the url, either find or create the
// corresponding page structure.  If the URL is deep, we also
// find the corresponding site structure.
func PageFromURL(ustr string) (pi *Page) {
	var ex bool
	if pi, ex = pages[ustr]; !ex {
		pi = &Page{
			URL:        ustr,
			Links:      make(map[string]*Page),
			Ignored:    make(map[string]int),
			CrawlState: NotCrawled,
		}
		pages[ustr] = pi
	}
	return pi
}

func (pm Pagemap) Get(url string) (p *Page) {
	if p, e := pages[url]; e {
		return p
	}
	return nil
}

func (pm Pagemap) Exists(url string) bool {
	if p := pm.Get(url); p != nil {
		return true
	}
	return false
}

func (pm Pagemap) Set(url string, p *Page) {
	pages[url] = p
}
