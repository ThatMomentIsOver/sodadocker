# sodadocker

简易的 Docker 漏洞扫描工具。利用 Docker Remote API 分析层级关系。

支持对指定的 docker 镜像内的 package 进行已知漏洞的检测。给定 Docker image ID 即可开始扫描。

## 其他安全风险检查

- 可检查 Docker 是否对外开放 ssh
