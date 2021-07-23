# hero-data-sdk
本框架是基于数数科技进行改进的。
支持配置2个http地址和2个日志目录，请按需配置。

## 1. 安装 SDK

`go get github.com/zhanqixuan/hero-data-sdk`

## 2.更新SDK

`go get -u  github.com/zhanqixuan/hero-data-sdk`

# 应用

## 该SDK支持2种方式，第一种是直接http请求推送，如下：

```
import "github.com/zhanqixuan/hero-data-sdk/herodata"
config := herodata.BatchConfig{
		ServerUrl:     "http://127.0.0.1:8089/api/sync/index", //必填
		AppId:         "test",
		ShuShuServerUrl: "", //推送数数科技，可选
		ShuShuAppId:     "",
		AutoFlush:     true,
		BatchSize:     100,
		Interval:      5,
		Compress:false,
	}
consumer, err := herodata.NewBatchConsumerWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

// 创建 TDAnalytics
ta := herodata.New(consumer)
// 设置事件属性
properties := map[string]interface{}{
    // 系统预置属性, 可选. "#time" 属性是系统预置属性，传入 time.Time 对象，表示事件发生的时间
    // 如果不填入该属性，则默认使用系统当前时间
	"#time": time.Now().UTC(),
    // 系统预置属性, 可选. 如果服务端中能获取用户 IP 地址, 并填入该属性
    // 数数会自动根据 IP 地址解析用户的省份、城市信息
    "#ip": "123.123.123.123",
    // 用户自定义属性, 字符串类型
    "prop_string": "abcdefg",
    // 用户自定义属性, 数值类型
    "prop_num": 56.56,
    // 用户自定义属性, bool 类型
    "prop_bool": true,
    // 用户自定义属性, time.Time 类型
	"prop_date": time.Now(),
}

account_id := "user_account_id" // 账号 ID
distinct_id := "ABCDEF123456AGDCDD" // 访客 ID

// 上报事件名为 TEST_EVENT 的事件. account_id 和 distinct_id 不能同时为空
ta.Track(account_id, distinct_id, "TEST_EVENT", properties)
//立即提交数据到相应的接收端
ta.Flush()

//关闭SDK
defer ta.Close()
```

## 第二种是先往本地写文件，然后通过flume插件往服务器kafka写入

```
// 创建按小时切分的 log consumer, 日志文件存放在当前目录
	// 创建按天切分的 log consumer, 不设置单个日志上限
	config := herodata.LogConfig{
		FileNamePrefix: "event",
		Directory:      "/var/log/hero_data",     //必填
		SecondDir:      "/var/log/hero_data_bak", //可选，备用日志地址，如果需要可以填写
	}
	consumer, err := herodata.NewLogConsumerWithConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ta := herodata.New(consumer)

	ta.SetSuperProperties(map[string]interface{}{
		"super_is_date":   time.Now(),
		"super_is_bool":   true,
		"super_is_string": "hello",
		"super_is_num":    15.6,
	})

	accountId := "AA"
	distinctId := "ABCDEF123456"
	properties := map[string]interface{}{
		// "#time" 属性是系统预置属性，传入 datetime 对象，表示事件发生的时间，如果不填入该属性，则默认使用系统当前时间
		//"#time":time.Now(),
		"update_time": time.Now(),
		// "#ip" 属性是系统预置属性，如果服务端中能获取用户 IP 地址，并填入该属性，数数会自动根据 IP 地址解析用户的省份、城市信息
		"#ip":     "123.123.123.123",
		"id":      "22",
		"catalog": "a",
		"is_boo":  true,
	}
	for i := 0; i < 100; i++ {

		// track事件
		err := ta.Track(accountId, distinctId, "view_page", properties)
		if err != nil {
			fmt.Println(err)
		}

	}
	ta.Flush()

	defer ta.Close()
  
```
