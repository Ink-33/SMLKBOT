# SMLKBOT
SMLKBOT, 基于[CQHTTPAPI](https://github.com/richardchien/coolq-http-api)的一个聚合群聊娱乐机器人  

## 1,功能
以下功能均已实装，并可以通过修改配置文件进行独立开关  


### 1.1, BiliAuCard
**注意：本功能需要CQP**  
提取聊天中的au号,然后返回音频分享卡片。  
效果：  
![au9](docs/au9.png)  

*注:不显示图片是TIM特性,图片均能在移动端QQ正常显示*

### 1.2, VTBMusic
**注意：本功能需要CQP**  
VTBMusic功能可以快捷地将您喜欢的歌曲分享给大家。所有音乐资源均来自于VTBMusic，请确认您要分享的歌曲已在VTBMusic正常上架。  

注意: 指令超时时间为**60秒**  

指令列表(不需要空格)：
```  
1: 普通点歌 -> vtb点歌+歌曲名
```
效果：  
![vtb点歌](docs/vtb1.png) 
```  
2: 精准点歌 -> vtbid点歌+歌曲ID  
```
效果：  
![vtbid点歌](docs/vtb2.png) 
```
3: 歌手点歌 -> vtb歌手+歌手名  
```
效果：  
![vtb歌手](docs/vtb3.png) 
```
4: 获取帮助 -> vtbhelp
```
效果：  
![vtbhelp](docs/vtbhelp.png)  

*注:不显示图片是TIM特性,图片均能在移动端QQ正常显示*

## 2,配置
将`conf.example.json`重命名为`conf.json`,
```json
{
    "CoolQ": {
        "Api": {
            "": {
                "HTTPAPIAddr": "",
                "HTTPAPIToken": "",
                "HTTPAPIPostSecret": ""
            }
        },
        "HTTPServer": {
            "ListeningPath": "/api/cqmsg",
            "ListeningPort": 12345
        }
    },
    "Feature": [
        {
            "BiliAu2Card": true,
            "VTBMusic": true
        }
    ]
}
```