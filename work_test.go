package hamtask

import (
	"fmt"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	w := NewWorker(3, func(i int, data Data) {
		fmt.Printf("no%d:%s\n", i, data.String())
	})
	w.Start() //启动消费者，开始等待生产者生产
	//生产者同步通道
	done := make(chan bool, 2)
	//生产者1
	go func() {
		for i := 0; i < 20; i++ {
			w.Put(Data{"hi"}) //生产
		}
		done <- true
	}()
	//生产者2
	go func() {
		for i := 0; i < 20; i++ {
			w.Put(String("hello")) //生产
		}
		done <- true
	}()

	//等待2个生产者生产完毕
	for i := 0; i < 2; i++ {
		<-done
	}

	w.Close() //结束生产队列，通知所有消费者

	w.Wait() //等待消费者执行完任务
}

func TestWorker1(t *testing.T) {
	w := NewWorker(3, func(i int, data Data) {
		fmt.Printf("no%d:%d\n", i, data.Int())
	})
	w.Start() //启动消费者，开始等待生产者生产

	w.Puts([]Data{Data{1}, Data{2}, Data{3}, Data{4}})

	w.Close() //结束生产队列，通知所有消费者

	w.Wait() //等待消费者执行完任务
}

func TestFullWorker(t *testing.T) {
	NewFullWorker(3, func(i int, data Data) {
		fmt.Printf("no%d:%d\n", i, data.Int())
	}, func(array chan Data) {
		for _, i := range []Data{Data{1}, Data{2}, Data{3}, Data{4}} {
			array <- i
		}
	}).Start() //启动消费者，开始等待生产者生产
}

/*
这里可以对比将n设置为1~100的数字，可以明显看出多线程的时间优势：
1个线程需要10秒
2个线程需要5秒
10个线程只需要1秒
*/
func TestSimpleWorker(t *testing.T) {
	NewSimpleWorker(1000, func(i int, data Data) {
		time.Sleep(time.Millisecond * 100)
		println(i, data.String())
	}, func() Data {
		return Data{"hi"}
	}, 1000).Start()
}

func TestData(t *testing.T) {
	d := Data{value: int64(12)}
	println(d.Type())
}
