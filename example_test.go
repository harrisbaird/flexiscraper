package flexiscraper_test

import (
	"os"
	"testing"

	. "github.com/harrisbaird/flexiscraper"
	. "github.com/harrisbaird/flexiscraper/q"
	"github.com/nbio/st"
)

type QwerteeItem struct {
	Name       string
	User       string
	ImageURL   string
	Prices     map[string]int
	LastChance bool
}

func TestParse(t *testing.T) {
	c := loadFixture("testdata/qwertee.html")
	var items []QwerteeItem

	c.Loop("//div[contains(@class, \"big-slide\") and contains(@class, \"tee\")]/div", func(i int, c *Context) {
		item := QwerteeItem{}
		item.Name = c.Find("@data-name")
		item.User = c.Find("@data-user")
		item.ImageURL = c.Build(XPath("@data-id"), Replace("https://www.qwertee.com/images/designs/zoom/%s.jpg")).String()
		item.LastChance = i > 2

		currencies := []string{"usd", "gbp", "eur"}
		prices := map[string]int{}
		for _, currency := range currencies {
			prices[currency] = c.Build(XPath("@data-tee-price-"+currency), Replace("%s00")).Int()
		}
		item.Prices = prices

		items = append(items, item)
	})

	st.Assert(t, len(items), 6)

	st.Assert(t, items[0].Name, "Black tee")
	st.Assert(t, items[0].User, "BlancaVidal")
	st.Assert(t, items[0].ImageURL, "https://www.qwertee.com/images/designs/zoom/111078.jpg")
	st.Assert(t, items[0].Prices, map[string]int{"gbp": 900, "eur": 1100, "usd": 1200})

	// Last 3 designs should be last chance
	lastChance := []bool{}
	for _, item := range items {
		lastChance = append(lastChance, item.LastChance)
	}
	st.Assert(t, lastChance, []bool{false, false, false, true, true, true})
}

func loadFixture(path string) *Context {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	c, err := Parse(file)
	if err != nil {
		panic(err)
	}

	return c
}
