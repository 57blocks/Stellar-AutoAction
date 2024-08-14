package constant

type HostEndPoint string

const (
	Host HostEndPoint = "http://st3llar-alb-365211.us-east-2.elb.amazonaws.com"
)

func (h HostEndPoint) String() string {
	return string(h)
}
