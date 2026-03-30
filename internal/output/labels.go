package output

var statusLabels = map[string]string{
	"settled":   "결제완료",
	"unsettled": "미결제",
	"draft":     "임시저장",
	"applying":  "신청중",
	"approved":  "승인",
	"rejected":  "반려",
	"cancelled": "취소",
	"sent":      "발송완료",
	"overdue":   "연체",
	"remanded":  "반송",
	"sending":   "발송중",
}

// StatusLabel returns the Korean label for a status value.
// Returns the original value if no translation exists.
func StatusLabel(status string) string {
	if label, ok := statusLabels[status]; ok {
		return label
	}
	return status
}
