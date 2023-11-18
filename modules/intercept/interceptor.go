package intercept

import (
	"apgo/system"
	"sync"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Interceptor struct {
	sync.RWMutex
	Sess *system.Session

	requestBinding binding.String
	requestPanel   *widget.Entry

	responsePanel *widget.Entry

	targetBinding binding.String
	targetEntry   *widget.Entry

	isIntercepting bool

	forwardChan chan bool
	dropChan    chan bool

	interceptSessionId int64
}

func NewInterceptor(s *system.Session) *Interceptor {
	p := &Interceptor{
		Sess:           s,
		requestBinding: binding.NewString(),
		targetBinding:  binding.NewString(),
		isIntercepting: false,
		forwardChan:    make(chan bool),
		dropChan:       make(chan bool),
	}
	return p
}
