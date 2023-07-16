# RBAC

## Quick Start

```shell
kubectl create namespace rbac
kubectl create -f svc-account.yaml -n rbac
kubectl create -f role-binding.yaml -n rbac
kubectl create -f role.yaml -n rbac
```

### More

Get ServiceAccount details

```shell
kubectl get sa -n rbac -o yaml
```

## References

* [26 | 基于角色的权限控制：RBAC](https://time.geekbang.org/column/article/42154)
