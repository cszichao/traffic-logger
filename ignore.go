package traffic_logger

type Ignore interface {
	Req(apiName string) bool
	Resp(apiName string) bool
}

type DefaultIgnore struct{}

func (DefaultIgnore) Req(apiName string) bool {
	return false
}

func (DefaultIgnore) Resp(apiName string) bool {
	return false
}

type IgnoreAll struct{}

func (IgnoreAll) Req(apiName string) bool {
	return true
}

func (IgnoreAll) Resp(apiName string) bool {
	return true
}
