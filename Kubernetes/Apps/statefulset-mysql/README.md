# StatefulSet MySQL

## Quick Start

```shell
kubectl create namespace statefulset-mysql
kubectl apply -f mysql-configmap.yaml -n statefulset-mysql
kubectl apply -f mysql-services.yaml -n statefulset-mysql
kubectl apply -f mysql-statefulset.yaml -n statefulset-mysql
```


### More

Write data to the MySQL by master
```shell
kubectl run mysql-client --image=mysql:5.7 -n statefulset-mysql -i --rm --restart=Never --\
 mysql -h mysql-0.mysql.statefulset-mysql <<EOF
CREATE DATABASE test;
CREATE TABLE test.messages (message VARCHAR(250));
INSERT INTO test.messages VALUES ('hello');
EOF
```

Read data from `mysql-read` service
```shell
kubectl run mysql-client --image=mysql:5.7 -n statefulset-mysql -i -t --rm --restart=Never --\
 mysql -h mysql-read.statefulset-mysql -e "SELECT * FROM test.messages"
```


Scaling the number of replicas

```shell
kubectl scale statefulset mysql -n statefulset-mysql --replicas=5
```

## References

* [20 | 深入理解StatefulSet（三）：有状态应用实践](https://time.geekbang.org/column/article/41366)

* [Kubernetes | Run a Replicated Stateful Application](https://kubernetes.io/docs/tasks/run-application/run-replicated-stateful-application/#statefulset)