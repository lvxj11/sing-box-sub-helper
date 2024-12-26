# sing-box-sub-helper
## sing-box订阅助手
功能非常单一，获取订阅节点列表，然后添加到模板生成sing-box的配置文件  
使用通用订阅。  
只覆写sing-box的outbound节点。对其他配置无影响。  
主要用于已经手搓了配置模板，又懒得搭建订阅转换后端的情况。

## 使用方法
1. 从Releases下载sing-box-sub-helper合适的执行文件
2. 运行sing-box-sub-helper，第一次运行会自动生成配置文件settings.ini
3. 编辑配置文件settings.ini，修改订阅地址、节点模板、输出文件路径等参数
4. 订阅地址和模板文件是必须的，输出文件可以使用默认
5. 再次运行sing-box-sub-helper，生成sing-box配置文件。
6. 运行sing-box run -c 配置文件。

### 其他说明
1. 暂时只支持一个订阅链接
2. 目前只适配了shadowsocks和trojan节点。（我只有这两种节点）
3. 为啥生成那么多临时文件？因为刚开始是用sh脚本写的，后来改用golang，为了方便调试就维持原状。
4. 为啥用golang？因为单文件部署无依赖，方便。
