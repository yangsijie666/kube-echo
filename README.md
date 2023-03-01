# kube-echo

## 简介

定义了一个命名空间维度的资源类型 Echo CR，在其 `spec.saySomeThing` 中指定字符串，其 Controller 会将该字段回写至 `status.echoResult` 中。

## 使用手册

### 部署

```shell
make deploy IMG="fightingsj/kube-echo:v1"
```

> `fightingsj/kube-echo` 为公共仓库，当然也可以替换其为自有 repository
> 若使用私有仓库，请在 k8s 的 `kube-echo-system` 中定义 Secret 并修改 `kube-echo-controller-manager` 这个 Deployment 中的 `imagePullSecret` 字段

> 在 `kube-echo-system` 中可以看到部署好的 Controller

### 卸载

```shell
make undeploy
```

### 测试用例

#### 新建

新建一个新的 Echo 实例：

```yaml
apiVersion: gogogo.yangsijie666.github.com/v1
kind: Echo
metadata:
  name: echo-sample
spec:
  saySomeThing: "testEcho-Hahaha"
```

在 k8s 中使用以下命令验证结果是否符合预期：

```shell
>> kubectl get echo
NAME          ECHORESULT
echo-sample   testEcho-Hahaha
```

或者也可以通过观察其 `status.echoResult` 看是否与 `spec.saySomeThing` 一致

```shell
>> kubectl get echo echo-sample -oyaml
apiVersion: gogogo.yangsijie666.github.com/v1
kind: Echo
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"gogogo.yangsijie666.github.com/v1","kind":"Echo","metadata":{"annotations":{},"name":"echo-sample","namespace":"default"},"spec":{"saySomeThing":"testEcho-Hahaha"}}
  creationTimestamp: "2023-03-01T05:22:37Z"
  generation: 1
  name: echo-sample
  namespace: default
  resourceVersion: "1681973681"
  selfLink: /apis/gogogo.yangsijie666.github.com/v1/namespaces/default/echoes/echo-sample
  uid: 46f1e2f1-30e6-4c29-b1be-efce4896b097
spec:
  saySomeThing: testEcho-Hahaha
status:
  echoResult: testEcho-Hahaha
  observedGeneration: 1
```

#### 更新

通过更新上述创建的 echo 实例中的 `spec.saySomeThing` 字段，观察是否 `status.echoResult` 最终与其保持一致

```shell
>> kubectl edit echo echo-sample
echo.gogogo.yangsijie666.github.com/echo-sample edited

>> kubectl get echo echo-sample
NAME          ECHORESULT
echo-sample   hoho

>> kubectl get echo echo-sample -oyaml
apiVersion: gogogo.yangsijie666.github.com/v1
kind: Echo
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"gogogo.yangsijie666.github.com/v1","kind":"Echo","metadata":{"annotations":{},"name":"echo-sample","namespace":"default"},"spec":{"saySomeThing":"testEcho-Hahaha"}}
  creationTimestamp: "2023-03-01T05:22:37Z"
  generation: 2
  name: echo-sample
  namespace: default
  resourceVersion: "1681988010"
  selfLink: /apis/gogogo.yangsijie666.github.com/v1/namespaces/default/echoes/echo-sample
  uid: 46f1e2f1-30e6-4c29-b1be-efce4896b097
spec:
  saySomeThing: hoho
status:
  echoResult: hoho
  observedGeneration: 2
```

> 从结果可以看出，saySomeThing 被改成 `hoho` 并成功更新至 `status` 中


#### 删除

```shell
>> kubectl delete echo echo-sample
echo.gogogo.yangsijie666.github.com "echo-sample" deleted

>> kubectl get echo 
No resources found in default namespace.
```

> 从结果可以看出，`echo-sample` 被成功删除
