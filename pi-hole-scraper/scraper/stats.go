package scraper

type Stats struct {
	AdsBlockedToday     int     `json:"ads_blocked_today"`
	AdsPercentageToday  float64 `json:"ads_percentage_today"`
	DomainsBeingBlocked int     `json:"domains_being_blocked"`
	DNSQueriesToday     int     `json:"dns_queries_today"`
}
