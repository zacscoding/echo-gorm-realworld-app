package types

type Status string

var (
	StatusDeleted = Status("deleted")
)

// StatusResponse represents a status response.
type StatusResponse struct {
	Status   Status                 `json:"status"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToStatusResponse converts given status string and metadata map to StatusResponse
func ToStatusResponse(status Status, meta map[string]interface{}) *StatusResponse {
	return &StatusResponse{
		Status:   status,
		Metadata: meta,
	}
}
