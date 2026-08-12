package main

type BusinessLine struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName,omitempty"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Discount    float64 `json:"discount,omitempty"`
	Amount      float64 `json:"amount"`
}

type Customer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Country     string `json:"country"`
	City        string `json:"city,omitempty"`
	Contact     string `json:"contact,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	Website     string `json:"website,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Rating      string `json:"rating"`
	Source      string `json:"source,omitempty"`
	Description string `json:"description,omitempty"`
	Orders      int    `json:"orders"`
	Total       string `json:"total"`
	Updated     string `json:"updated"`
	Archived    bool   `json:"archived,omitempty"`
}

type CustomerCommunication struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerId"`
	Type       string `json:"type"`
	Subject    string `json:"subject"`
	Content    string `json:"content"`
	Contact    string `json:"contact,omitempty"`
	OccurredAt string `json:"occurredAt"`
	CreatedAt  string `json:"createdAt"`
	CreatedBy  string `json:"createdBy"`
}

type Quote struct {
	ID       string         `json:"id"`
	Subject  string         `json:"subject"`
	Customer string         `json:"customer"`
	Value    string         `json:"value"`
	Valid    string         `json:"valid"`
	Status   string         `json:"status"`
	Products string         `json:"products"`
	Owner    string         `json:"owner,omitempty"`
	Currency string         `json:"currency,omitempty"`
	Freight  float64        `json:"freight,omitempty"`
	Tax      float64        `json:"tax,omitempty"`
	Terms    string         `json:"terms,omitempty"`
	Lines    []BusinessLine `json:"lines"`
	Updated  string         `json:"updated"`
}

type Order struct {
	ID       string         `json:"id"`
	Customer string         `json:"customer"`
	Quote    string         `json:"quote,omitempty"`
	PO       string         `json:"po,omitempty"`
	Products string         `json:"products"`
	Value    string         `json:"value"`
	Currency string         `json:"currency,omitempty"`
	Status   string         `json:"status"`
	Delivery string         `json:"delivery"`
	Terms    string         `json:"terms,omitempty"`
	Progress int            `json:"progress"`
	Lines    []BusinessLine `json:"lines"`
	Updated  string         `json:"updated"`
}

type Document struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Customer string         `json:"customer"`
	Order    string         `json:"order"`
	Template string         `json:"template"`
	Language string         `json:"language,omitempty"`
	Status   string         `json:"status"`
	Value    string         `json:"value,omitempty"`
	Lines    []BusinessLine `json:"lines"`
	Updated  string         `json:"updated"`
}

type Product struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Price        string `json:"price"`
	Stock        int    `json:"stock"`
	HS           string `json:"hs"`
	Status       string `json:"status"`
	Description  string `json:"desc,omitempty"`
	Tags         string `json:"tags,omitempty"`
	Updated      string `json:"updated,omitempty"`
}

type Supplier struct {
	ID            string `json:"id"`
	Company       string `json:"company"`
	Phone         string `json:"phone,omitempty"`
	Contact       string `json:"contact,omitempty"`
	Email         string `json:"email,omitempty"`
	Product       string `json:"product,omitempty"`
	Specification string `json:"specification,omitempty"`
	Quote         string `json:"quote,omitempty"`
	Notes         string `json:"notes,omitempty"`
	Source        string `json:"source,omitempty"`
	Updated       string `json:"updated"`
	Archived      bool   `json:"archived,omitempty"`
}

type PrivateFile struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Size     string   `json:"size"`
	Bytes    int64    `json:"bytes,omitempty"`
	SHA256   string   `json:"sha256,omitempty"`
	Mime     string   `json:"mime,omitempty"`
	Source   string   `json:"source"`
	Status   string   `json:"status"`
	Path     string   `json:"-"`
	Updated  string   `json:"updated"`
}

type AgentMessageRecord struct {
	ID         string   `json:"id"`
	AgentID    string   `json:"agentId"`
	SessionKey string   `json:"sessionKey"`
	Role       string   `json:"role"`
	Text       string   `json:"text"`
	Sources    []string `json:"sources,omitempty"`
	RunID      string   `json:"runId,omitempty"`
	CreatedAt  string   `json:"createdAt"`
}

type NewsItem struct {
	ID        string `json:"id"`
	Category  string `json:"category"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	SourceURL string `json:"sourceUrl,omitempty"`
	Time      string `json:"time"`
	Relevance string `json:"relevance"`
}

type Recommendation struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Time   string `json:"time"`
}

type Job struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Schedule  string `json:"schedule"`
	Enabled   bool   `json:"enabled"`
	BuiltIn   bool   `json:"builtIn"`
	LastRun   string `json:"lastRun,omitempty"`
	NextRun   string `json:"nextRun,omitempty"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	UpdatedAt string `json:"updatedAt"`
}

type Plugin struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Enabled      bool     `json:"enabled"`
	Status       string   `json:"status"`
	Capabilities []string `json:"capabilities"`
	UpdatedAt    string   `json:"updatedAt"`
}

type UserPreferences struct {
	RecommendationEnabled bool                `json:"recommendationEnabled"`
	NewsShowLimit         int                 `json:"newsShowLimit"`
	NewsFrequency         string              `json:"newsFrequency"`
	NewsCountries         string              `json:"newsCountries"`
	NewsTopics            string              `json:"newsTopics"`
	NewsSources           string              `json:"newsSources"`
	AgentAllowlists       map[string][]string `json:"agentAllowlists"`
}
