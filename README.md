# Coolpy7 Community

[![Release](https://img.shields.io/github/v/release/Coolpy7/Coolpy7_Community)](https://github.com/Coolpy7/Coolpy7_Community/releases)
[![WebSite](https://img.shields.io/website?up_message=Coolpy7&url=http%3A%2F%2Fwww.coolpy.net%2F)](http://www.coolpy.net/)
[![License](https://img.shields.io/github/license/Coolpy7/Coolpy7_Community)](https://github.com/Coolpy7/Coolpy7_Community/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Coolpy7/Coolpy7_Community)](https://goreportcard.com/report/github.com/Coolpy7/Coolpy7_Community)
[![OpenIssue](https://img.shields.io/github/issues/Coolpy7/Coolpy7_Community)](https://github.com/Coolpy7/Coolpy7_Community/issues)
[![ClosedIssue](https://img.shields.io/github/issues-closed/Coolpy7/Coolpy7_Community)](https://github.com/Coolpy7/Coolpy7_Community/issues?q=is%3Aissue+is%3Aclosed)
![Stars](https://img.shields.io/github/stars/Coolpy7/Coolpy7_Community)

## 简介

Coolpy7社区版 是基于Epoll为通信核心开发的高性能MQTT服务库，专注于 物联网\即时通信\消息总线等 场景。

## 特性

* **已经支持** 
    - `MQTT协议支持` 支持MQTT3.1/MQTT3.1.1
    - `Qos全支持` 支持Qos0 Qos1 Qos2
    - `WebSocket` 提供高性能Websocket桥接MQTT服务
    - `JWT验证` 提供JWT密钥验证（token于password连接参数）
    - 支持 `Linux`（操作系统）
    - 防空连接攻击
    - 防`Self Ddos`攻击

# 快速开始

### 1. 下载

[这里](https://github.com/Coolpy7/Coolpy7_Community/releases) 下载最新版本的各系统平台编译运行文件。

### 2 服务参数

```
$ ./coolpy7_community_linux -h
$ Usage of ./coolpy7_community_linux:
  -j string
        jwt secret key(multiple split by ,
  -k string
        tls key file path
  -l string
        host port (default 1883) (default ":1883")
  -p string
        tls pem file path
  -w string
        host ws port (default 8083) (default ":8083")
  -wk string
        wss key file path
  -wp string
        wss pem file path
```

参数说明:

* `-l` tcp服务端口参数。 默认值1883端口，例示 :1883。
* `-w` websocket服务端口参数。 默认值8083端口，例示 :8083。
* `-j` jwt密钥参数。 默认值为空时MQTT服务不做任何身份验证，多个密钥是以半角豆号（,)分隔。例示 jwtsk1,jwtsk2,jwtsk3。
* `-p` tls服务pem证书文件绝对路径。 
* `-k` tls服务key证书文件绝对路径。
* `-wp` websocket wss服务pem证书文件绝对路径。
* `-wk` websocket wss服务key证书文件绝对路径。

### 3 运行服务

参数说明:
* 启动服务前请先确认相关执行权限

```
$ ./coolpy7_community_xxx -l=:1883 -w=:8083
$ 2023/06/21 17:57:16 Coolpy7 Community On Port :1883
$ 2023/06/21 17:57:16 Coolpy7 Community Websocket On Port :8083

```

### 4 安全关闭服务

* kill -2 pid

## 应用示例

* **多种技术客户端示例** [APP(Android,Flutter)--后端(Golang,Java,Nodejs,Python3,C#)--前端(Electron,React,Vue,微信小程序,WebSocket)--单片机(ESP8266)](https://github.com/Coolpy7/mqtt-client-examples)

- **Web浏览器应用示例**，[这里](https://github.com/Coolpy7/mqtt_web_browser_client)

- **微信小程序聊天室示例**，[这里](https://github.com/Coolpy7/wxsmallapp)

- **Web浏览器聊天室示例（可与微信小程序示例连同一Coolpy7后互相聊天）**：[这里](https://github.com/Coolpy7/Cp7Chat)

- **Web浏览器mqtt.js客户端示例**， [这里](https://github.com/Coolpy7/mqttjs_browser_client_demo)

## 性能

### Qos0压力测试
```
./mqtt-stresser-linux-amd64 -broker tcp://192.168.0.7:7006 -num-clients 100 -num-messages 1500 -rampup-delay 1s -rampup-size 10 -global-timeout 30s -timeout 10s -constant-payload 100
```
```
# Configuration
Concurrent Clients: 100
Messages / Client:  150000

# Results
Published Messages: 150000 (100%)
Received Messages:  108852 (73%)
Completed:          6 (6%)
Errors:             1 (1%)
- ConnectFailed:      0 (0%)
- SubscribeFailed:    0 (0%)
- TimeoutExceeded:    1 (100%)
- Aborted:            93 (93%)

# Publishing Throughput
Fastest: 218333 msg/sec
Slowest: 18942 msg/sec
Median: 73253 msg/sec

  < 38881 msg/sec  17%
  < 58820 msg/sec  36%
  < 78759 msg/sec  56%
  < 98698 msg/sec  69%
  < 118637 msg/sec  75%
  < 138577 msg/sec  80%
  < 158516 msg/sec  85%
  < 178455 msg/sec  90%
  < 198394 msg/sec  94%
  < 218333 msg/sec  99%
  < 238272 msg/sec  100%

# Receiving Througput
Fastest: 132334 msg/sec
Slowest: 10472 msg/sec
Median: 19992 msg/sec

  < 22658 msg/sec  59%
  < 34844 msg/sec  93%
  < 47030 msg/sec  96%
  < 59216 msg/sec  97%
  < 83589 msg/sec  98%
  < 95775 msg/sec  99%
  < 144520 msg/sec  100%
```

### Qos1压力测试
```
./mqtt-stresser-linux-amd64 -broker tcp://192.168.0.7:7006 -num-clients 100 -num-messages 1500 -rampup-delay 1s -rampup-size 10 -global-timeout 30s -timeout 10s -constant-payload 100 -publisher-qos 1  -subscriber-qos 1
```
```
# Configuration
Concurrent Clients: 100
Messages / Client:  150000

# Results
Published Messages: 150000 (100%)
Received Messages:  150000 (100%)
Completed:          100 (100%)
Errors:             0 (0%)

# Publishing Throughput
Fastest: 7524 msg/sec
Slowest: 4690 msg/sec
Median: 6028 msg/sec

  < 4974 msg/sec  7%
  < 5257 msg/sec  16%
  < 5540 msg/sec  33%
  < 5824 msg/sec  39%
  < 6107 msg/sec  52%
  < 6391 msg/sec  63%
  < 6674 msg/sec  67%
  < 6957 msg/sec  74%
  < 7241 msg/sec  86%
  < 7524 msg/sec  99%
  < 7807 msg/sec  100%

# Receiving Througput
Fastest: 141490 msg/sec
Slowest: 24047 msg/sec
Median: 51706 msg/sec

  < 35792 msg/sec  8%
  < 47536 msg/sec  41%
  < 59280 msg/sec  68%
  < 71025 msg/sec  80%
  < 82769 msg/sec  88%
  < 94513 msg/sec  92%
  < 106257 msg/sec  93%
  < 118002 msg/sec  94%
  < 129746 msg/sec  95%
  < 141490 msg/sec  99%
  < 153234 msg/sec  100%
```

### Qos2压力测试
```
./mqtt-stresser-linux-amd64 -broker tcp://192.168.0.7:7006 -num-clients 100 -num-messages 1500 -rampup-delay 1s -rampup-size 10 -global-timeout 30s -timeout 10s -constant-payload 100 -publisher-qos 2  -subscriber-qos 2
```
```
# Configuration
Concurrent Clients: 100
Messages / Client:  150000

# Results
Published Messages: 150000 (100%)
Received Messages:  150000 (100%)
Completed:          100 (100%)
Errors:             0 (0%)

# Publishing Throughput
Fastest: 4493 msg/sec
Slowest: 2737 msg/sec
Median: 3380 msg/sec

  < 2913 msg/sec  14%
  < 3088 msg/sec  33%
  < 3264 msg/sec  41%
  < 3440 msg/sec  53%
  < 3615 msg/sec  61%
  < 3791 msg/sec  70%
  < 3967 msg/sec  79%
  < 4142 msg/sec  86%
  < 4318 msg/sec  97%
  < 4493 msg/sec  99%
  < 4669 msg/sec  100%

# Receiving Througput
Fastest: 145966 msg/sec
Slowest: 38577 msg/sec
Median: 74970 msg/sec

  < 49316 msg/sec  7%
  < 60055 msg/sec  20%
  < 70794 msg/sec  38%
  < 81533 msg/sec  59%
  < 92272 msg/sec  71%
  < 103010 msg/sec  80%
  < 113749 msg/sec  84%
  < 124488 msg/sec  90%
  < 135227 msg/sec  96%
  < 145966 msg/sec  99%
  < 156705 msg/sec  100%
```
