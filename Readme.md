## 带屏推广告设备

#### MQTT Config

##### Broker：

待定



##### Server <------> client: 

​	server: Publish topic

​	client: Subscribe topic





##### client status report

​	client: Publish topic



### Communication Protocol



#### server to client



##### 询问实时在线设备

服务器通过此指令获得实时在线的设备信息，中间层会发送实时在线的设备信息（通常用于服务器上线后第一次初始化）

```json
{
    "Typ": "askDevices"
}
```

> - `Typ` 指令类型



##### 初始化设备信息

用于回应中间层的`askDevices` 指令

```json
{
    "Typ": "init",
    "Devices": [
        {
            "Typ": "UIT",
            "Mac": "12345678"
        },
        
        {
            "Typ": "UIT",
            "Mac": "12345678"
        }
    ]
}
```

> - `Typ` 指令类型，`init`表示服务要告诉中间层有哪些设备需要管理和上报信息，一般用于回应中间层`askDevices`指令，添加或删除设备时，直接发送新的`init`指令即可
> - `Devs` 设备对象数组
>   - `Typ` 设备类型
>   - `Mac` 设备表示符，一般为Mac地址
>   - `WriteProp` 设备写属性对象
>     - `ImgAd` 广告图片数组



##### 发送状态信息

设置设备的属性

```json
{
    "Typ": "status",
    "Devices": [
        {
            "Typ": "UIT",
            "Mac": "12345678",
            "WriteProp": {
                "ImgAd": [
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5"
                ]
            }
        },
        {
            "Typ": "UIT",
            "Mac": "12345678",
            "WriteProp": {
                "ImgAd": [
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5",
                    "http://path/md5"
                ]
            }
        }
    ]
}
```

> - `Typ` 指令类型
> - `Devices` 设备对象数组
>   - `Typ` 设备类型
>   - `Mac` 设备表示符，一般为Mac地址
>   - `WriteProp`设备属性对象
>     - `ImgAd` 广告图片数组



##### 指定设备版本信息

```json
{
    "Typ": "version",
	"DevTyp":"UIT",
    "DevRom":"http://path/ota.bin"
}
```

> `Typ` 指令类型
>
> `DevTyp` 设备类型
>
> `DevRom` 设备rom下载地址



#### client to server



##### 询问需要管理的设备

中间层发送此指令来询问服务器需要管理的设备（每次中间层启动后发送）

```json
{
    "Typ": "askDevices"
}
```

> - `Typ` 指令类型



##### 实时在线设备上报

用于回应服务器的`askDevices`指令

```json
{
    "Typ": "init",
    "Devices": [
		{
            "Typ": "UIT",
            "Mac": "12345678",
    	},
        
		{
            "Typ": "UIT",
            "Mac": "123124124",
    	}  
    ]
}
```



> - `Typ` 指令类型`
> - `Device ` 设备对象数组
> - `Typ` 设备类型
> - `Mac` 设备标志符





#### client report

##### 设备状态上报

单设备单信息上报

```json
{
    "Typ": "status",
    "Device": {
        "Typ": "UIT",
        "Mac": "12345678",
        "Status": "offline",
        "ReadProp": {
            "ImgProgres":0
        }
    }
}
```



> - `Typ` 指令类型
> - `Device`设备状态对象
> - `Typ` 设备类型
> - `Mac` 设备标志符
> - `Status` 设备在线状态
>   - `online` 设备在线
>   - `offline` 设备离线
> - `ReadProp` 设备只读属性对象
>   - `ImgProgress` 图片下发进度，数值类型
>     - `0` 未发布，
>     - `1` 发布中
>     - `2` 发布完成



