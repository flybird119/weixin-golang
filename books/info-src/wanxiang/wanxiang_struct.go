package wanxiang

type wanxiang struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Msg    string `json:"msg"`
		Status string `json:"status"`
		Result book   `json:"result"`
	} `json:"result"`
}

type book struct {
	Summary   string `json:"summary"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Edition   string `json:"edition"`
	Pubdate   string `json:"pubdate"`
	Pic       string `json:"pic"`
	Publisher string `json:"publisher"`
	Price     string `json:"price"`
	Isbn      string `json:"isbn"`
	Isbn10    string `json:"isbn10"`
}
