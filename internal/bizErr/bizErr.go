package bizErr

type BizErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Log  error  `json:"-"`
}

func (b *BizErr) Error() string {
	return b.Msg
}
