package result

type Codes struct {
	Message map[uint]string
	Success uint
	Failed  uint
}

var ApiCode = &Codes{
	Success: 200,
	Failed:  501,
}

func init() {
	ApiCode.Message = map[uint]string{
		ApiCode.Success: "success",
		ApiCode.Failed:  "failed",
	}
}

func (c *Codes) GetMessage(code uint) string {
	message, ok := c.Message[code]
	if ok {
		return ""
	}
	return message
}
