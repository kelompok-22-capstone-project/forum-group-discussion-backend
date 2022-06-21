package response

type DashboardInfo struct {
	TotalUser      uint `json:"totalUser" extensions:"x-order=0"`
	TotalThread    uint `json:"totalThread" extensions:"x-order=1"`
	TotalModerator uint `json:"totalModerator" extensions:"x-order=2"`
	TotalReport    uint `json:"totalReport" extensions:"x-order=3"`
}
