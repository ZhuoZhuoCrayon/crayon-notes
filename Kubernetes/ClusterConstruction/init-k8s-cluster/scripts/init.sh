## 环境配置
# 更新软件源
yum update -y
# 内核参数配置
TOTAL_MEM=$(free -b | awk 'NR==2{print $2}')
TOTAL_MEM=${TOTAL_MEM:-$(( 16 * 1024 * 1024 *1024 ))}
PAGE_SIZE=$(getconf PAGE_SIZE)
PAGE_SIZE=${PAGE_SIZE:-4096}
THREAD_SIZE=$(( PAGE_SIZE << 2 ))
sed -ri.k8s.bak '/k8s config begin/,/k8s config end/d' /etc/sysctl.conf
cat >> "/etc/sysctl.conf" << EOF
# 为了支持 k8s service, 必须开启
net.ipv4.ip_forward=1
# 使得 Linux 节点的 iptables 能够正确查看桥接流量
net.bridge.bridge-nf-call-ip6tables=1
net.bridge.bridge-nf-call-iptables=1
# k8s config end
EOF
sysctl --system
cat > /etc/security/limits.d/99-k8s.conf << EOF
# k8s config begin
*   soft  nproc  1028546
*   hard  nproc  1028546
*   soft  nofile  204800
*   hard  nofile  204800
# k8s config end
EOF

# 设置主机名
hostnamectl set-hostname 127-0-0-3-node
# 添加域名解析
sed -ri.k8s.hosts.bak '/k8s hosts begin/,/k8s hosts end/d' /etc/hosts
cat >> "/etc/hosts" << EOF
# k8s hosts begin
127.0.0.1 k8scp
127.0.0.1 127-0-0-1-master
127.0.0.2 127-0-0-2-node
127.0.0.3 127-0-0-3-node
127.0.0.4 127-0-0-4-nfs-node
# k8s hosts end
EOF

# 禁用主机交换分区，如果不满足，系统会有一定几率出现 io 飙升
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
# 禁用 selinux
sed -i 's/enforcing/disabled/' /etc/selinux/config
setenforce 0
cat /etc/selinux/config

# 重启确保配置生效
reboot

## 初始化 Docker
# 安装 yum 工具
yum install -y -q yum-utils

# 创建 docker repo
MIRROR_URL="http://mirrors.tencentyun.com"
YUM_REPO="${MIRROR_URL}/docker-ce/linux/centos/docker-ce.repo"
curl -fs "$YUM_REPO" | sed "s#https://download.docker.com#${MIRROR_URL}/docker-ce#g" | tee "$HOME/docker-ce.repo"
# 若主机OS为 TencentOS 需要修改 repo
sed -i "s/\$releasever/7/g" "$HOME/docker-ce.repo"
# 导入 docker repo
yum-config-manager --add-repo "$HOME/docker-ce.repo"
yum makecache

# 安装 containerd.io
# 用于解决 `yum -y install docker-ce-cli-19.03.15-3.el7 docker-ce-19.03.15-3.el7 containerd.io`
# 报错：requires containerd.io >= 1.2.2-3, but none of the providers can be installed
# 参考：https://mebee.info/2021/04/13/post-32969/
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/containerd.io-1.2.2-3.3.el7.x86_64.rpm
dnf install containerd.io-1.2.2-3.3.el7.x86_64.rpm

# 安装 docker
# el7 代表什么？
# EL 是Red Hat E nterprise L inux（EL）的缩写，EL7 是 Red Hat 7.x，CentOS 7.x 和CloudLinux 7.x 的下载，需根据实际情况填写。
yum -y install docker-ce-cli-19.03.15-3.el7 docker-ce-19.03.15-3.el7

# 创建依赖目录
# 在 data 目录下创建 docker 的数据目录
install -dv /data/lib/docker
# 在 etc 目录下创建 docker 配置文件目录
install -dv /etc/docker

# 创建 docker 配置文件
cat << EOF | sudo tee /etc/docker/daemon.json
{
  "data-root": "/data/lib/docker",
  "exec-opts": ["native.cgroupdriver=systemd"],
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com"
  ],
  "max-concurrent-downloads": 10,
  "live-restore": false,
  "log-level": "info",
  "log-opts": {
    "max-size": "100m",
    "max-file": "5"
  },
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ]
}
EOF

# 在 Docker 启动时将 FORWARD chain 默认策略设为 ACCEPT
install -dv /etc/systemd/system/docker.service.d/
cat << EOF | sudo tee /etc/systemd/system/docker.service.d/k8s-docker.conf
[Service]
ExecStartPost=/sbin/iptables -P FORWARD ACCEPT
EOF
systemctl daemon-reload
# 托管 Docker 进程
systemctl enable --now docker

# 检查 docker.service 配置
# 会有以下2个配置文件
# /usr/lib/systemd/system/docker.service
# /etc/systemd/system/docker.service.d/k8s-docker.conf
systemctl cat docker

# 验证 Docker 功能可用
docker run --rm docker.io/library/hello-world:latest
docker run -it ubuntu bash

## 安装 Kubeadm

# 添加源
cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kubelet kubeadm kubectl
EOF

# 安装依赖
yum -y install kubectl-1.20.12-0 kubeadm-1.20.12-0 kubelet-1.20.12-0 --disableexcludes=kubernetes

# 检查工具版本
kubectl version --client --short
kubeadm version -o short

## 启动 kubelet
install -dv /data/lib/kubelet
cat << EOF | sudo tee /etc/sysconfig/kubelet
KUBELET_EXTRA_ARGS="--root-dir=/data/lib/kubelet"
EOF
systemctl enable --now kubelet


# 创建 etcd 数据目录
install -dv /data/lib/etcd

# 创建 Kubeadm 配置
cat << EOF | tee "$HOME/kubeadm-config.yaml"
apiVersion: kubeadm.k8s.io/v1beta2
apiServer:
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
# 使用Node alias而不是IP
controlPlaneEndpoint: k8scp:6443
controllerManager: {}
dns:
  type: CoreDNS
etcd:
  local:
    # etcd数据目录推荐放到数据盘中
    dataDir: /data/lib/etcd
imageRepository: k8s.gcr.io
kind: ClusterConfiguration
kubernetesVersion: v1.20.12
networking:
  dnsDomain: cluster.local
  podSubnet: 10.244.0.0/16
  serviceSubnet: 10.96.0.0/12
scheduler: {}
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
nodeRegistration:
  name: 127-0-0-1-master
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
# kube-proxy 采用 ipvs 模式
mode: ipvs
EOF
