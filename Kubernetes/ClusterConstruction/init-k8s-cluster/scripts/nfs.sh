# 查看磁盘信息
fdisk -l

# 创建文件系统格式
mkfs -t ext4 /dev/vdb

# 创建并将新建的文件分区挂载到 `/nfs-storage/`
mkdir -p /nfs-storage/data
mount /dev/vdb /nfs-storage

# 获取磁盘的 `UUID`
blkid /dev/vdb

# UUID=c6ffe18c-274e-4cec-a3b5-068baa8ce9de /nfs-storage ext4 defaults 0 0
vim /etc/fstab


# 检查 /etc/fstab 文件是否写入成功
mount -a
df -TH

# 检查依赖是否安装
rpm -qa | grep nfs
rpm -qa | grep rpcbind

# 安装并启动依赖
yum install rpcbind nfs-utils
systemctl start rpcbind
systemctl enable rpcbind
systemctl start nfs-server
systemctl enable nfs-server

## NFS 共享目录文件配置
## 参考：https://www.huweihuang.com/linux-notes/tools/nfs-usage.html
sed -ri.k8s.bak '/k8s config begin/,/k8s config end/d' /etc/exports
cat >> "/etc/exports" << EOF
# k8s config begin
/nfs-storage/data 10-0-4-9-master(rw,insecure,sync,no_subtree_check,no_root_squash)
/nfs-storage/data 127-0-0-4-nfs-node(rw,insecure,sync,no_subtree_check,no_root_squash)
/nfs-storage/data 10-0-8-6-node(rw,insecure,sync,no_subtree_check,no_root_squash)
/nfs-storage/data 10-0-8-9-node(rw,insecure,sync,no_subtree_check,no_root_squash)
# k8s config end
EOF

# -a 全部挂载或者全部卸载
# -r 重新挂载
exportfs -ra
showmount -e
fuser -m -v /nfs-storage

ss | grep 2049
# 查看其他节点挂 NFS 的情况
netstat | grep :nfs

## 搭建 NFS 客户端
# 安装依赖
yum install -y nfs-utils

# 挂载
mkdir -p /nfs-storage/data
mount -t nfs 127-0-0-4-nfs-node:/nfs-storage/data /nfs-storage/data

# 开机自动挂载
sed -ri.k8s.bak '/k8s config begin/,/k8s config end/d' /etc/fstab
cat >> "/etc/fstab" << EOF
# k8s config begin
127-0-0-4-nfs-node:/nfs-storage/data    /nfs-storage/data    nfs    defaults    0 0
# k8s config end
EOF

mount -a

df -h
