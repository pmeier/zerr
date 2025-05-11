package zerr

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

type Zerr struct {
	HTTPStatusCode int
	Redacted       bool
	e              *zerolog.Event
	msg            string
	b              *bytes.Buffer
	text           string
}

func New(msg string, opts ...OptFunc) *Zerr {
	o := resolveOptions(opts)
	var b bytes.Buffer
	l := zerolog.New(&b)
	return &Zerr{HTTPStatusCode: o.httpStatusCode, Redacted: o.redacted, e: l.Log(), msg: msg, b: &b}
}

func (z *Zerr) Send() *Zerr {
	if z.text == "" {
		z.text = z.b.String()
	}

	return z
}

func (z *Zerr) Msg(msg string) *Zerr {
	z.msg = msg
	return z.Send()
}

func (z *Zerr) MsgFunc(createMsg func() string) *Zerr {
	return z.Msg(createMsg())
}

func (z *Zerr) Msgf(format string, v ...any) *Zerr {
	return z.Msg(fmt.Sprintf(format, v...))
}

func (z *Zerr) Error() string {
	z.Send()
	return z.text
}

func (z *Zerr) MarshalZerologObject(e *zerolog.Event) {
	var m map[string]any
	if err := json.Unmarshal([]byte(z.Error()), &m); err != nil {
		panic(err.Error())
	}

	if msg, ok := m[zerolog.MessageFieldName]; ok {
		delete(m, zerolog.MessageFieldName)
		e.Str(zerolog.ErrorFieldName, msg.(string))
	}

	for k, v := range m {
		e.Any(k, v)
	}
}
