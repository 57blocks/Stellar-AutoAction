package model

type Lambda struct {
	ICUD
}

func (l *Lambda) TableName() string {
	return "object_lambda"
}

func (l *Lambda) TableNameWithAbbr() string {
	return "object_lambda AS ol"
}
