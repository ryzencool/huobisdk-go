package robot

// AccountsResponse huobi accounts
type AccountsResponse struct {
	Status string `json:"status"`
	Data   []struct {
		ID      int    `json:"id"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		State   string `json:"state"`
	} `json:"data"`
}

// OrderResponse huobi place order
type OrderResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// DepthResponse huobi get market depth
type DepthResponse struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	Tick   struct {
		ID   int64       `json:"id"`
		Ts   int64       `json:"ts"`
		Bids [][]float64 `json:"bids"`
		Asks [][]float64 `json:"asks"`
	} `json:"tick"`
}

// OrderDetailResponse huobi get order detail
type OrderDetailResponse struct {
	Status string `json:"status"`
	Data   struct {
		ID              int    `json:"id"`
		Symbol          string `json:"symbol"`
		AccountID       int    `json:"account-id"`
		Amount          string `json:"amount"`
		Price           string `json:"price"`
		CreatedAt       int64  `json:"created-at"`
		Type            string `json:"type"`
		FieldAmount     string `json:"field-amount"`
		FieldCashAmount string `json:"field-cash-amount"`
		FieldFees       string `json:"field-fees"`
		FinishedAt      int64  `json:"finished-at"`
		UserID          int    `json:"user-id"`
		Source          string `json:"source"`
		State           string `json:"state"`
		CanceledAt      int    `json:"canceled-at"`
		Exchange        string `json:"exchange"`
		Batch           string `json:"batch"`
	} `json:"data"`
}

// MatchConfig huobi get config
type MatchConfig struct {
	AccessID  string  `json:"accessID"`
	SecretKey string  `json:"secretKey"`
	Host      string  `json:"host"`
	Match     []Match `json:"match"`
	Place     []Place `json:"place"`
}

// Match match
type Match struct {
	Symbol   string `json:"symbol"`
	Limit    int    `json:"limit"`
	Strategy struct {
		Day struct {
			Normal []struct {
				AmountMin int `json:"amountMin"`
				AmountMax int `json:"amountMax"`
				Interval  int `json:"interval"`
			} `json:"normal"`
			Special []struct {
				AmountMin int `json:"amountMin"`
				AmountMax int `json:"amountMax"`
				Interval  int `json:"interval"`
			} `json:"special"`
		} `json:"day"`
		Night struct {
			Normal []struct {
				AmountMin int `json:"amountMin"`
				AmountMax int `json:"amountMax"`
				Interval  int `json:"interval"`
			} `json:"normal"`
			Special []struct {
				AmountMin int `json:"amountMin"`
				AmountMax int `json:"amountMax"`
				Interval  int `json:"interval"`
			} `json:"special"`
		} `json:"night"`
	} `json:"strategy"`
}

// Place config
type Place struct {
	Symbol string `json:"symbol"`
	Step   []Step `json:"step"`
}

// Step config
type Step struct {
	FlagStart int     `json:"flagStart"`
	FlagEnd   int     `json:"flagEnd"`
	AmountMin float64 `json:"amountMin"`
	AmountMax float64 `json:"amountMax"`
}

// CurOrdersResponse current orderresponse
type CurOrdersResponse struct {
	Status string `json:"status"`
	Data   []struct {
		ID              int    `json:"id"`
		Symbol          string `json:"symbol"`
		AccountID       int    `json:"account-id"`
		Amount          string `json:"amount"`
		Price           string `json:"price"`
		CreatedAt       int64  `json:"created-at"`
		Type            string `json:"type"`
		FieldAmount     string `json:"field-amount"`
		FieldCashAmount string `json:"field-cash-amount"`
		FieldFees       string `json:"field-fees"`
		FinishedAt      int64  `json:"finished-at"`
		UserID          int    `json:"user-id"`
		Source          string `json:"source"`
		State           string `json:"state"`
		CanceledAt      int    `json:"canceled-at"`
		Exchange        string `json:"exchange"`
		Batch           string `json:"batch"`
	} `json:"data"`
}
