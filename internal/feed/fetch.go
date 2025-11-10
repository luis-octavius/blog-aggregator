package feed 

import (
	"net/http"
	"context"
	"fmt"
	"encoding/xml"
	"io"
	"html"

	"github.com/luis-octavius/blog-aggregator/internal/types"
)

// FetchFeed retrieves and parses an RSS feed from the specified URL. 
// it createas an HTTP request with context and custom User-Agent header,
// then unmarshals the XML response into an RSSFeed struct.
// the function also handles HTML unescaping for text field 

// returns a parsed RSS feed or an error if any of these steps fails
// HTTP request creation or execution; response body reading; XML unmarshaling
// HTML unescaping
func FetchFeed(ctx context.Context, feedURL string) (*types.RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// set custom User-Agent to indentify this app
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var rssFeed types.RSSFeed 
	if err = xml.Unmarshal(body, &rssFeed); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// clean HTML from text fields 
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	// process each item to unescape HTML content 
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title) 
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item 
	}	

	return &rssFeed, nil
}
