package douban

type DoubanBook struct {
	Msg       string   `json:"msg"`
	Subtitle  string   `json:"subtitle"`
	Author    []string `json:"author"`
	Publisher string   `json:"publisher"`
	Isbn10    string   `json:"isbn10"`
	Isbn13    string   `json:"isbn13"`
	Title     string   `json:"title"`
	Price     string   `json:"price"`
	Summary   string   `json:"summary"`
	Images    struct {
		Large string `json:"large"`
	}
	Author_intro string `json:"author_intro"`
	Pubdate      string `json:"pubdate"`
}
