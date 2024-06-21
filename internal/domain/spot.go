package domain

type Spot struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Status  string `json:"status"`
    EventID int    `json:"event_id"`
}
