# BiliAu2Card  
基于[CQHTTPAPI](https://github.com/richardchien/coolq-http-api),提取聊天中的au号,然后返回音频卡片。  
注意：卡片分享功能需要酷Q Pro, Air版将以文字形式发送。

## 用法

将`conf.example.json`重命名为`cong.json`  

```shell
$mv ./conf.example.json ./cong.json
```
### 配置文件范本
```json
{
    "CoolQ": [
        {
            "Api": {
                "HttpAPIAddr": "",
                "HttpAPIToken": "",
            }
        }
    ],
    "BiliAu2Card": [    
        {
            "ListeningPath": "/api/cqmsg",
            "ListeningPort": 65122
        }
    ]
}
```
### 参数说明:  
`HttpAPIAddr`: 此处填写你的CQHTTPAPI的http请求地址  
`HttpAPIToken`: 此处填写你的CQHTTPAPI的http请求Token  

`BiliAu2Card`: 将此处填写到CQHTTPAPI的上报地址中  

`ListeningPath`: BiliAu2Card监听PATH  
`ListeningPort`: BiliAu2Card监听端口

修改完成后运行即可。

## 将来将要加入的功能
1: Add HMAC-SHA1 signature verification.