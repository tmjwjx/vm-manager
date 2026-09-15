# vm-manager

基于 Go 的 **Proxmox VE 虚拟机管理服务**,提供虚拟机从创建到销毁的全生命周期管理。

> 2025-03 ~ 2025-04 团队项目(2 人),本人为主要开发者(约 3/4 提交)。

## 功能

- **虚拟机全生命周期**:创建(自动分配最小可用 VMID)、启动、重启、销毁、状态与信息查询
- **WebSocket 实时通道**:前端通过 WebSocket 下发操作指令,异步回调回传执行结果
- **网络配置**:SSH 登录虚拟机设置静态 IP(含 CIDR → 子网掩码换算)、自动获取虚机 IP
- **台账管理**:MySQL 持久化虚拟机信息,按集群/宿主机维度管理

## 架构

```
vm-manager/
├── pve/               # Proxmox VE 对接服务(主体)
│   ├── cmd/           # 入口
│   ├── controller/    # WebSocket 控制器(StartVM / RebootVM / CreateVM / DestroyVM ...)
│   ├── logic/         # 业务逻辑(异步回调、HTTP 客户端)
│   ├── models/        # 请求/响应结构
│   └── pkg/           # globals · utils
└── vm/                # 虚拟机台账模块(Gin + GORM + MySQL + Viper)
```

## 技术栈

- Go 1.22+
- GORM(gorm v1 / v2 双版本)+ MySQL
- Viper 配置管理
- golang.org/x/crypto(SSH)
- WebSocket
- Proxmox VE API / qm 命令行

## 说明

- 项目开发于校内环境,原托管于内网 GitLab,此为归档镜像
