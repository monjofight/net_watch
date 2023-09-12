package netflix

type Title struct {
	ID        int
	Name      string
	Seasons   []Season
	CreatedAt string
	UpdatedAt string
}

func NewTitle(TitleID int) *Title {
	return &Title{
		ID: TitleID,
	}
}

type Season struct {
	ID        int
	TitleID   int
	Name      string
	Episodes  []Episode
	CreatedAt string
	UpdatedAt string
}

type Episode struct {
	ID        int
	TitleID   int
	SeasonID  int
	Name      string
	Image     string
	Watched   bool
	CreatedAt string
	UpdatedAt string
}

type SearchQuery struct {
	Query   string
	Results []SearchResult
}

type SearchResult struct {
	ID      int
	Title   string
	Link    string
	Snippet string
}

func NewSearchQuery(query string) *SearchQuery {
	return &SearchQuery{
		Query: query,
	}
}
