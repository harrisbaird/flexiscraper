package flexiscraper_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	. "github.com/harrisbaird/flexiscraper"
	. "github.com/harrisbaird/flexiscraper/q"
	"github.com/nbio/st"
)

type QwerteeItem struct {
	Name           string
	User           string
	ImageURL       string
	OtherImageURLs []string
	Prices         map[string]int
	LastChance     bool
}

func TestExample(t *testing.T) {
	r, err := recorder.New("testdata/qwertee")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop()

	scraper := New()
	scraper.HTTPClient = &http.Client{Transport: r}

	domain := scraper.NewDomain("http://qwertee.com")
	c, err := domain.FetchRoot()
	if err != nil {
		panic(err)
	}

	var items []QwerteeItem

	c.Each("//div[contains(@class, \"big-slide\") and contains(@class, \"tee\")]/div", func(i int, c *Context) {
		item := QwerteeItem{}
		item.Name = c.Attr("data-name")
		item.User = c.Attr("data-user")
		item.ImageURL = c.Build(Attr("data-id"), Replace("https://www.qwertee.com/images/designs/zoom/%s.jpg")).String()
		item.OtherImageURLs = c.Build(XPath(".//source/@srcset"), Replace("https://www.qwertee.com%s")).StringSlice()
		item.LastChance = i > 2

		currencies := []string{"usd", "gbp", "eur"}
		prices := map[string]int{}
		for _, currency := range currencies {
			prices[currency] = c.Build(Attr("data-tee-price-"+currency), Replace("%s00")).Int()
		}
		item.Prices = prices

		items = append(items, item)
	})

	st.Assert(t, len(items), 6)

	st.Assert(t, items[0].Name, "Kawaii Cats")
	st.Assert(t, items[0].User, "jozephine")
	st.Assert(t, items[0].ImageURL, "https://www.qwertee.com/images/designs/zoom/109681.jpg")
	st.Assert(t, items[0].Prices, map[string]int{"gbp": 900, "eur": 1100, "usd": 1200})

	// Last 3 designs should be last chance
	lastChance := []bool{}
	for _, item := range items {
		lastChance = append(lastChance, item.LastChance)
	}
	st.Assert(t, lastChance, []bool{false, false, false, true, true, true})
}
