package main

import (
	"fmt"
	"hero_data_sdk/herodata"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	// 创建 BatchConsumer, 指定接收端地址、APP ID、上报批次
	config := herodata.BatchConfig{
		ServerUrl:     "http://127.0.0.1:8088/", //必填
		AppId:         "test",
		ShuShuServerUrl: "", //推送数数科技，可选
		ShuShuAppId:     "",
		AutoFlush:     true,
		BatchSize:     100,
		Interval:      5,
	}
	consumer, err := herodata.NewBatchConsumerWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建 TDAnalytics
	ta := herodata.New(consumer)

	// 设置公共事件属性
	ta.SetSuperProperties(map[string]interface{}{
		"super_string": "value",
		"super_bool":   false,
	})

	fmt.Printf("%v", time.Now())
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(threadName string) {
			defer wg.Done()

			account_id := threadName
			distinct_id := "ABCDEF123456"
			properties := map[string]interface{}{
				// "#time" 属性是系统预置属性，传入 datetime 对象，表示事件发生的时间，如果不填入该属性，则默认使用系统当前时间
				"time_now": time.Now(),
				// "#ip" 属性是系统预置属性，如果服务端中能获取用户 IP 地址，并填入该属性，数数会自动根据 IP 地址解析用户的省份、城市信息
				"#ip": "123.123.123.123",
				"id":  "12",
				//#uuid  去重，服务端比较稳定，可不填，如果填，按照以下标准8-4-4-4-12的String()
				"#uuid":   "6ba7b810-9dad-2233-80b4-00c04fd430c8",
				"catalog": "p",
				"bool":    true,
				"aa":      12,
			}
			// track事件
			err := ta.Track(account_id, distinct_id, "ViewProduct", properties)
			if err != nil {
				fmt.Println(err)
			}
		}(fmt.Sprintf("thread%d", i))
	}
	fmt.Printf("%v", time.Now())

	time.Sleep(1000 * time.Second)

	wg.Wait()
	ta.Flush()
	defer ta.Close()
}
