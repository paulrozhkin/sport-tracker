package dto

type Healthcheck struct {
	TotalDbInvokes      int64 `json:"totalDbInvokes"`
	CurrentDbConnection int   `json:"currentDbConnection"`
	MaxDbConnections    int   `json:"maxDbConnections"`
}
