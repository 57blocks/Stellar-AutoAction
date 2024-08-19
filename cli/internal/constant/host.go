package constant

type EndPoint string

const (
	// Host TODO: confirmation when deployed
	Host EndPoint = "http://st3llar-alb-365211.us-east-2.elb.amazonaws.com"
)

func (h EndPoint) String() string {
	return string(h)
}
