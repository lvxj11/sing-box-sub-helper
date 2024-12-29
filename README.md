# sing-box-sub-helper
## sing-box订阅助手
功能非常单一，获取订阅节点列表，然后添加到模板生成sing-box的配置文件  
使用通用订阅。  
只覆写sing-box的outbound节点。对其他配置无影响。  
主要用于已经手搓了配置模板，又懒得搭建订阅转换后端的情况。

## 使用方法
1. 从Releases下载sing-box-sub-helper合适的执行文件。
2. 运行`sing-box-sub-helper`，第一次运行会自动生成配置文件settings.ini。
3. 编辑配置文件settings.ini，修改订阅地址、节点模板、输出文件路径等参数。
4. subscribeURL，templatePath 和 outputPath 是必须的。分别是订阅地址，模板文件和和输出文件路径。
5. 本程序只能解析通用订阅，用浏览器访问订阅地址应该能获取到base64编码的密文内容。如果不是将不能正确解析。
6. 如果修改了开始步骤，则根据步骤需要配置相应的数据源文件路径。
7. 再次运行sing-box-sub-helper，生成sing-box配置文件。
8. 运行`sing-box check -c "配置文件"`检查配置文件，或`sing-box run -c "配置文件"`直接执行。
9. 后台运行或添加到系统服务方法请自行google。

### 其他说明
1. 暂时只支持一个订阅链接
2. 已适配节点协议：
    - shadowsocks
    - vmess
    - trojan
    - hysteria2
3. 为啥生成那么多临时文件？因为刚开始是用sh脚本写的，改用golang后为了方便调试就维持原状。后来发现可以从不同步骤开始，以免频繁获取订阅被限制访问。
4. 为啥用golang？因为单文件部署无依赖，方便。
