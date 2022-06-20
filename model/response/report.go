package response

type Report struct {
	ID                string `json:"ID" extensions:"x-order=0"`
	ModeratorID       string `json:"moderatorID" extensions:"x-order=1"`
	ModeratorUsername string `json:"moderatorUsername" extensions:"x-order=2"`
	ModeratorName     string `json:"moderatorName" extensions:"x-order=3"`
	UserID            string `json:"userID" extensions:"x-order=4"`
	Username          string `json:"username" extensions:"x-order=5"`
	Name              string `json:"name" extensions:"x-order=6"`
	Reason            string `json:"reason" extensions:"x-order=7"`
	Status            string `json:"status" extensions:"x-order=8"`
	// ReportedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	ReportedOn string `json:"reportedOn" extensions:"x-order=9"`
}
