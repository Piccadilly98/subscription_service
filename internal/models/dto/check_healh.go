package dto

type CheckHealth struct {
	ServerStatus string  `json:"server_status"`
	DBStatus     string  `json:"db_status"`
	Error        *string `json:"error_text,omitempty"`
}

func ToCheckHealthDTO(serverStatus, dbStatus string, err error) *CheckHealth {
	dto := &CheckHealth{
		ServerStatus: serverStatus,
		DBStatus:     dbStatus,
	}

	if err != nil {
		errStr := err.Error()
		dto.Error = &errStr
	}

	return dto
}
