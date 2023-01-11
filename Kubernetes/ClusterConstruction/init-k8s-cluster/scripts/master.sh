## master 节点配置

# kubeadm 初始化
kubeadm init --config=kubeadm-config.yaml --upload-certs | tee kubeadm-init.out

# kube config 配置
export KUBECONFIG=/etc/kubernetes/admin.conf
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

# 生成节点加入集群的命令
kubeadm token create --print-join-command

# 网络配置
kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml

# 通过部署 Nginx 快速体验
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --port 80 --type=NodePort


# 安装 Helm
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

# NFS provisioner
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
helm install nfs-subdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner --set nfs.server=10-0-4-15-nfs-node --set nfs.path=/nfs-storage/data/ --set storageClass.defaultClass=true

# test
helm repo add https://zhuozhuocrayon.github.io/helm-charts/
helm repo add myrepo  https://zhuozhuocrayon.github.io/helm-charts/
helm search repo djangocli --versions
helm install djangocli myrepo/djangocli --version=0.5.8 -f values-private.yaml

# k9s
wget https://github.com/derailed/k9s/releases/download/v0.25.18/k9s_Linux_x86_64.tar.gz
tar -zxf k9s_Linux_x86_64.tar.gz -C /usr/local/bin
k9s help
k9s info

# 性能指标
wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
# L140 加上配置
#         - --kubelet-insecure-tls=true
#         - --kubelet-preferred-address-types=InternalIP
vim components.yaml
kubectl apply -f components.yaml
