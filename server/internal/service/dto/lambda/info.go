package lambda

import "github.com/57blocks/auto-action/server/internal/model"

type (
	ReqInfo struct {
		Lambda string `uri:"lambda"`
	}

	RespInfo struct {
		ID           uint64      `json:"id"`
		VpcID        uint64      `json:"vpc_id"`
		FunctionName string      `json:"function_name"`
		FunctionArn  string      `json:"function_arn"`
		Runtime      string      `json:"runtime"`
		Role         string      `json:"role"`
		Handler      string      `json:"handler"`
		Description  string      `json:"description"`
		CodeSHA256   string      `json:"code_sha256"`
		Version      string      `json:"version"`
		RevisionID   string      `json:"revision_id"`
		VPC          RespVpc     `json:"vpc" gorm:"foreignKey:vpc_id"`
		Schedulers   []Scheduler `json:"schedulers" gorm:"foreignKey:lambda_id"`
	}

	RespVpc struct {
		ID               uint64        `json:"-" gorm:"column:id"`
		AmazonID         string        `json:"vpc_id" gorm:"column:aws_id"`
		SubnetIDs        model.StrList `json:"subnet_ids" gorm:"type:text[]"`
		SecurityGroupIDs model.StrList `json:"security_group_ids" gorm:"type:text[]"`
	}

	Scheduler struct {
		LambdaID    uint64 `json:"lambda_id"`
		ScheduleArn string `json:"schedule_arn"`
		Expression  string `json:"expression"`
	}
)
