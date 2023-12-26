# CloudTrailCollector 数据采集器
目前支持数据源
输入：  
国内腾讯云审计日志、华为云审计日志
海外AWS云审计日志、阿里云审计日志
输出：
ES

# 原理
配置分成input和output，一个input和output构成一个采集管道，支持配置多个input和output，管道数N*M关系  
output使用condition字段来匹配寻找对应的input
  
目前支持input类型有
1. timer 每2s中输出一个字符串，用于测试。
2. tc_audit_meta_data 腾讯云审计日志
3. hw_audit_meta_data 华为云审计日志
4. aws_cloud_trail AWS云审计日志
5. ali_cloud_trail 阿里云审计日志

目前支持output类型有
1. elastic ES输出
2. console 数据打印到终端供测试使用

# 扩展
input编写
```
GetName() string
该接口用于获取input名称，打印使用
Update(context context.Context, ch chan interface{})
该接口用于从获取数据，ch是输出数据管道

使用registerInputCollector注册
```
output编写
```
GetName() string
该接口用于获取output名称，打印使用
Insert(context context.Context, body interface{}) error
该接口写入body数据

使用registerOutputCollector注册
```
# 配置
``` 
// 每个输出管道的工作并发数， 提高insert并发
worker = 4

//input 数组
[[input]]
type = "aws_cloud_trail"
endpoint = ""
access_key = "xxxxx" #ak
secret_key = "xxxxx" #sk
region_id = "xxxxx"

[[output]]
type = "console"
condition = ""  // 所有input都会与之匹配
index = ""
end_points = []
username = ""
password = ""

[[output]]
type = "elastic"
condition = "type == aws_cloud_trail"
index = "aws_cloud_operation"
end_points = ["http://127.0.0.q:9200"]
username = "xxx"
password = "xxx"
```




 












