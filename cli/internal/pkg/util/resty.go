package util

type Error struct {
	Message string `json:"message"`
	Notes   string `json:"notes"`
	Status  int    `json:"status"`
}

func (e Error) Error() string {
	return e.Message
}
