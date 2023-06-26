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
./mqtt-stresser-linux-amd64 -broker tcp://172.18.157.151:1883 -num-clients 10 -num-messages 150 -rampup-delay 1s -rampup-size 10 -global-timeout 180s -timeout 20s
```
```
# Configuration
Concurrent Clients: 10
Messages / Client:  1500

# Results
Published Messages: 1500 (100%)
Received Messages:  992 (66%)
Completed:          0 (0%)
Errors:             0 (0%)

# Publishing Throughput
Fastest: 84124 msg/sec
Slowest: 19982 msg/sec
Median: 63101 msg/sec

  < 26396 msg/sec  20%
  < 45639 msg/sec  30%
  < 52053 msg/sec  40%
  < 58467 msg/sec  50%
  < 77710 msg/sec  80%
  < 84124 msg/sec  90%
  < 90538 msg/sec  100%

# Receiving Througput
Fastest: 109125 msg/sec
Slowest: 18135 msg/sec
Median: 28378 msg/sec

  < 27234 msg/sec  50%
  < 36333 msg/sec  60%
  < 45432 msg/sec  80%
  < 90927 msg/sec  90%
  < 118224 msg/sec  100%
```

### Qos1压力测试
```
./mqtt-stresser-linux-amd64 -broker tcp://172.18.157.151:1883 -num-clients 10 -num-messages 150 -rampup-delay 1s -rampup-size 10 -global-timeout 180s -timeout 20s -publisher-qos 1  -subscriber-qos 1
```
```
..........
# Configuration
Concurrent Clients: 10
Messages / Client:  1500

# Results
Published Messages: 1500 (100%)
Received Messages:  1500 (100%)
Completed:          10 (100%)
Errors:             0 (0%)

# Publishing Throughput
Fastest: 5699 msg/sec
Slowest: 5046 msg/sec
Median: 5295 msg/sec

  < 5111 msg/sec  10%
  < 5176 msg/sec  20%
  < 5241 msg/sec  30%
  < 5307 msg/sec  60%
  < 5437 msg/sec  80%
  < 5503 msg/sec  90%
  < 5764 msg/sec  100%

# Receiving Througput
Fastest: 59043 msg/sec
Slowest: 25028 msg/sec
Median: 31189 msg/sec

  < 28429 msg/sec  40%
  < 31831 msg/sec  50%
  < 35232 msg/sec  70%
  < 38634 msg/sec  80%
  < 52240 msg/sec  90%
  < 62445 msg/sec  100%
```

### Qos2压力测试
```
./mqtt-stresser-linux-amd64 -broker tcp://172.18.157.151:1883 -num-clients 10 -num-messages 150 -rampup-delay 1s -rampup-size 10 -global-timeout 180s -timeout 20s -publisher-qos 2  -subscriber-qos 2
```
```
# Configuration
Concurrent Clients: 10
Messages / Client:  1500

# Results
Published Messages: 1500 (100%)
Received Messages:  1500 (100%)
Completed:          10 (100%)
Errors:             0 (0%)

# Publishing Throughput
Fastest: 3187 msg/sec
Slowest: 2841 msg/sec
Median: 3101 msg/sec

  < 2876 msg/sec  10%
  < 2910 msg/sec  20%
  < 2945 msg/sec  30%
  < 3049 msg/sec  40%
  < 3083 msg/sec  50%
  < 3152 msg/sec  80%
  < 3187 msg/sec  90%
  < 3221 msg/sec  100%

# Receiving Througput
Fastest: 36544 msg/sec
Slowest: 12518 msg/sec
Median: 22943 msg/sec

  < 14921 msg/sec  20%
  < 17323 msg/sec  40%
  < 22128 msg/sec  50%
  < 26934 msg/sec  70%
  < 31739 msg/sec  80%
  < 34141 msg/sec  90%
  < 38946 msg/sec  100%
```
