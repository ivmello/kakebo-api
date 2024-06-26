package reports

type SummarizeOutput struct {
	Total           int    `json:"total"`
	TotalFormated   string `json:"total_formated"`
	Debits          int    `json:"debits"`
	DebitsFormated  string `json:"debits_formated"`
	Credits         int    `json:"credits"`
	CreditsFormated string `json:"credits_formated"`
}

type SummarizeByMonthOutput struct {
	Month           int    `json:"month"`
	Year            int    `json:"year"`
	Total           int    `json:"total"`
	TotalFormated   string `json:"total_formated"`
	Debits          int    `json:"debits"`
	DebitsFormated  string `json:"debits_formated"`
	Credits         int    `json:"credits"`
	CreditsFormated string `json:"credits_formated"`
}
