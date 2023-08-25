package xhttpres

type ResOne struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewResOne(d interface{}, msg string, codes ...int) ResOne {
	code := 200
	switch r := d.(type) {
	case bool:
		if !r {
			code = 500
		}
		d = nil
	}

	if len(codes) > 0 {
		code = codes[0]
	}

	if msg == "" {
		msg = getMsg(code)
	}

	return ResOne{
		Code: code,
		Msg:  msg,
		Data: d,
	}
}

func getMsg(code int) string {
	switch code {
	case 500:
		return "Fail"
	default:
		return "SUCCESS"
	}
}
