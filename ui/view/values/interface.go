package values

import (
	"gioui.org/io/event"
)

type View interface {
	OnLogin(body []byte)
	OnEvent(e event.Event, data interface{})
	Weight() WeightType
	ViewIndex() ViewIndexType
}
