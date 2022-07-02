package response

type Report struct {
	ID                string `json:"ID" extensions:"x-order=0"`
	ModeratorID       string `json:"moderatorID" extensions:"x-order=1"`
	ModeratorUsername string `json:"moderatorUsername" extensions:"x-order=10"`
	ModeratorName     string `json:"moderatorName" extensions:"x-order=11"`
	UserID            string `json:"userID" extensions:"x-order=12"`
	Username          string `json:"username" extensions:"x-order=13"`
	Name              string `json:"name" extensions:"x-order=14"`
	Reason            string `json:"reason" extensions:"x-order=15"`
	Status            string `json:"status" extensions:"x-order=16"`
	ThreadID          string `json:"threadID" extensions:"x-order=17"`
	ThreadTitle       string `json:"threadTitle" extensions:"x-order=18"`
	// ReportedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	ReportedOn string `json:"reportedOn" extensions:"x-order=19"`
	Comment    string `json:"comment" extensions:"x-order=20"`
	// ReportedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	CommentPublishedOn string `json:"commentPublishedOn" extensions:"x-order=21"`
}
