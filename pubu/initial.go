package pubu

import (
	"log"
	"sync"
	"time"
)

type askChannel chan ask

type ask struct {
	Type askType
	msg  string
}

type askType int

const (
	debug askType = iota
	warning
	mistake
	info
	good
	bad
)

type pubu struct {
	hook string
	ask  askChannel
}

var pb *pubu
var once sync.Once

func New(hook string) *pubu {
	once.Do(
		func() {
			ask := make(askChannel, 12)
			pb = &pubu{
				hook: hook,
				ask:  ask,
			}
			go start(ask)
			pb.Good("零信的初始化工作完成。")
		})

	return pb
}

func start(askChan askChannel) {
	beginTime := time.Now()
	for ask := range askChan {
		switch ask.Type {
		case get:
			data, err := ec.Get(ask.Path)
			ask.AnswerChan <- answer{body: data, err: err}
		case post:
			data, err := ec.Post(ask.Path, ask.Headers, ask.Body)
			ask.AnswerChan <- answer{body: data, err: err}
		default:
			log.Println("Wrong ask type.")
		}
		//由于pubu.im有API访问次数限制，“每个接入调用限制为每秒 10 次”
		//所以，需要暂停一下。
		gu.HoldOn(time.Millisecond*100, &beginTime)
	}
}
