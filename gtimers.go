package gtimers

import (
	"fmt"
	"sync"
	"time"
)

//定时器
type gtimer struct {
	Starttime   int64        `json:"starttime,omitempty"`   //开始时间
	Repeatcount int          `json:"repeatcount,omitempty"` //执行次数.0:重复执行,n>0,执行n次
	Interval    string       `json:"interval,omitempty"`    //任务周期
	Status      int          `json:"status,omitempty"`      //任务执行状态,0:未执行,1:执行中,2:执行完成
	timer       *time.Ticker //循环执行
}

type callback func()

func NewTimer() *gtimer {
	return &gtimer{
		Starttime:   time.Now().Unix(),
		Repeatcount: 0,
		Interval:    "1s", //默认定时周期一秒触发
		Status:      0,
	}
}

func (t *gtimer) SetStarttime(v int64) {
	t.Starttime = v
}

func (t *gtimer) SetRepeatcount(v int) {
	t.Repeatcount = v
}

func (t *gtimer) SetInterval(v string) {
	t.Interval = v
}

func (t *gtimer) SetStatus(v int) {
	t.Status = v
}

func (t *gtimer) Reset(v *gtimer) {
	t = v
}
func (t *gtimer) Stop() {
	t.timer.Stop()
}

func (t *gtimer) SyncStart(fn callback) error {
	var tongbu sync.WaitGroup
	nowtime := time.Now().Unix()
	interval, err := time.ParseDuration(t.Interval)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if t.Starttime > 0 && t.Starttime > nowtime {
		time.Sleep(time.Duration((t.Starttime - nowtime)) * time.Second) //开始时间-当前时间 = 段时间,阻塞一段时间
	}
	t.timer = time.NewTicker(interval)

	defer t.timer.Stop()
	if t.Repeatcount == 0 {
		for {
			select {
			case <-t.timer.C:
				tongbu.Add(1)
				go func() {
					fn()
					tongbu.Done()
				}()

			}

			tongbu.Wait()

		}
	} else {
		for {
			select {
			case <-t.timer.C:
				tongbu.Add(1)
				go func() {
					fn()
					tongbu.Done()
				}()

			}

			tongbu.Wait()
			t.Repeatcount--
			if t.Repeatcount <= 0 {
				//任务执行完成
				t.SetStatus(2)
				break
			}

		}
	}

	return nil
}

func (t *gtimer) Start(fn callback) error {

	nowtime := time.Now().Unix()
	interval, err := time.ParseDuration(t.Interval)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if t.Starttime > 0 && t.Starttime > nowtime {
		time.Sleep(time.Duration((t.Starttime - nowtime)) * time.Second) //开始时间-当前时间 = 段时间,阻塞一段时间
	}
	t.timer = time.NewTicker(interval)

	defer t.timer.Stop()
	if t.Repeatcount == 0 {
		for {
			select {
			case <-t.timer.C:
				go fn()
			}

		}
	} else {
		for {
			select {
			case <-t.timer.C:
				go fn()
			}

			t.Repeatcount--
			if t.Repeatcount <= 0 {
				//任务执行完成
				t.SetStatus(2)
				break
			}

		}
	}

	return nil
}

//func main() {
//	timer := time.NewTicker(15 * time.Second)
//	defer timer.Stop()
//	times := 0
//	//var tb sync.WaitGroup
//	for {
//		select {
//		case ss := <-timer.C:
//			fmt.Println(ss.Format("05"))
//			//tb.Add(1)
//			go func() {
//				time.Sleep(time.Second * 10)
//				fmt.Println("做一点事情")
//				//	tb.Done()
//			}()
//		}
//		//	tb.Wait()
//		times++
//		if times == 5 {
//			fmt.Println("任务执行完成...")
//			break
//		}
//	}

//}
