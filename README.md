NFS Client for Golang
=====================

Installation
------------

This library uses gb to build. Simply run gb in the root of the checkout then have a look in the {{_bin}} directory for examples.
=======
# nfs-proxy
# nfs-proxy 设计目的主要是用来做一个高可用的针对k8s的nfs下文件服务器
# 背景
* 在使用k8s过程中，随着容器增多小文件逐渐增多，达到1000w的时候，我们当初设计使用的Mfs无法承担如此流量
* 在容器的使用中，往往挂载着都是代码文件，这些文件往往是高读高写的，在这个场景中使用nfs能得到更好的使用速度
* 但是使用nfs的高可用方案特别少，而且通过k8s挂载nfs如果nfs挂掉，容器必须重新启动才能使用，有什么可以做到不重启容器呢，我们首先想到的是通过一层nfs-proxy作为代理进行管理nfs

# 目标
* 完成nfs-proxy代理层
* 完成nfs各个服务器检测
* 完成nfs副本集
* 完成nfs文件直接的平衡
* 快速部署nfs服务器，达到横向扩展
