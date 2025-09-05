package kb

import (
	"motadataAssignment/models"
	"sort"
	"strings"

	"github.com/google/uuid"
)

// KB holds articles
type KB struct {
	articles []models.Article
}

// New returns a KB pre-loaded with sample articles (hardcoded)
func New() *KB {
	k := &KB{
		articles: []models.Article{
			{
				ID:      uuid.NewString(),
				Title:   "How to reset your password",
				Content: "To reset your password go to account settings -> reset password. If you don't get an email check spam.",
			},
			{
				ID:      uuid.NewString(),
				Title:   "Troubleshooting network connectivity",
				Content: "Check that cable is plugged in, ensure DHCP is enabled, try 'ipconfig /renew' on Windows or 'dhclient' on Linux.",
			},
			{
				ID:      uuid.NewString(),
				Title:   "Installing Node.js on Ubuntu",
				Content: "Use apt to install node: curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash - && sudo apt-get install -y nodejs",
			},
			{
				ID:      uuid.NewString(),
				Title:   "Email SMTP misconfiguration",
				Content: "If sending fails, verify SMTP host, port, TLS settings and credentials. Use telnet to test connection.",
			},
		},
	}
	return k
}

// All returns all articles
func (k *KB) All() []models.Article {
	return k.articles
}

// FindRelevant returns articles ranked by simple keyword overlap score
func (k *KB) FindRelevant(query string, topN int) []models.Article {
	if query == "" {
		return []models.Article{}
	}
	qterms := tokenize(query)

	type scored struct {
		a     models.Article
		score int
	}
	var scoredList []scored
	for _, a := range k.articles {
		text := strings.ToLower(a.Title + " " + a.Content)
		score := 0
		for _, t := range qterms {
			if strings.Contains(text, t) {
				score++
			}
		}
		if score > 0 {
			scoredList = append(scoredList, scored{a: a, score: score})
		}
	}
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score > scoredList[j].score
	})

	res := []models.Article{}
	for i := 0; i < len(scoredList) && i < topN; i++ {
		res = append(res, scoredList[i].a)
	}
	return res
}

func tokenize(s string) []string {
	// naive tokenizer: lowercase, split on non-alphanumeric characters
	s = strings.ToLower(s)
	t := []rune{',', '.', ';', ':', '/', '\\', '?', '!', '-'}
	for _, r := range t {
		s = strings.ReplaceAll(s, string(r), " ")
	}
	parts := strings.Fields(s)
	return parts
}
