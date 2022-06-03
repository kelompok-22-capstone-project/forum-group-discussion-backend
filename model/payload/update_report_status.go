package payload

type UpdateReportStatus struct {
	// Status, available options: rejected, accepted
	Status string `json:"status" validate:"nonzero,min=5,max=9" extensions:"x-order=0"`
}
