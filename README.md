# Weplanx API

基于 Gin 设计的后端单体服务实验项目

## 快速部署

该方式是低代码项目后端的通用方案，用于无需定制的场景（预发布是测试版本）：

- ghcr.io/weplanx/api:latest
- ccr.ccs.tencentyun.com/weplanx/api:latest（国内）

案例将使用 Kubernetes 部署编排，复制部署内容（根据情况修改）：

1. 设置配置

```yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api.cfg
data:
  config.yml: |
    address: ":9000"
    # 不设置将无法获取到客户端 ip，也可以使用七层反向代理做上游（例如：全站加速、应用负载均衡等）
    trusted_proxies:
      - 10.42.0.0/16
    namespace: <应用名称>
    key: <32位密文>
    database:
      uri: mongodb://<username>:<password>@<host>:<port>/<database>?authSource=<authSource>
      dbName: <database>
    redis:
      uri: redis://<user>:<password>@<host>:<port>/<db_number>
    nats:
      hosts: [ ]
      nkey:
    cors:
      allowOrigins:
        - https://app.****.com
      allowMethods:
        - POST
      allowHeaders:
        - Content-Type
        - Accept
      allowCredentials: true
      maxAge: 7200
    passport:
      aud: [ 'console' ]
      exp: 720
    qcloud:
      secret_id: <secret_id>
      secret_key: <secret_key>
      cos:
        bucket: examplebucket-****
        region: ap-guangzhou
        expired: 300
    engines:
      pages:
        event: true
      users:
        projection:
          password: 0
```

2. 部署

```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: api
  name: api-deploy
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
        - image: ccr.ccs.tencentyun.com/weplanx/api:latest
          imagePullPolicy: Always
          name: api
          ports:
            - containerPort: 9000
          volumeMounts:
            - name: config
              mountPath: "/app/config"
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: api.cfg
            items:
              - key: "config.yml"
                path: "config.yml"
```

3. 设置入口

```yml
apiVersion: v1
kind: Service
metadata:
  name: api-svc
spec:
  ports:
    - port: 80
      protocol: TCP
      targetPort: 9000
  selector:
    app: api

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api
#   annotations:
#     traefik.ingress.kubernetes.io/router.tls: "true"
#     traefik.ingress.kubernetes.io/router.tls.certresolver: ****
#     traefik.ingress.kubernetes.io/router.tls.domains.0.main: ****.com
#     traefik.ingress.kubernetes.io/router.tls.domains.0.sans: "*.****.com"
spec:
  rules:
    - host: api.****.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-svc
                port:
                  number: 80
```

## 滚动更新

复制模板内容，并需要自行定制触发条件，原理是每次patch将模板中 `${tag}` 替换为版本执行

```yml
spec:
  template:
    spec:
      containers:
        - image: ccr.ccs.tencentyun.com/weplanx/api:${tag}
          name: api
```

例如：在 Github Actions
中 `patch deployment api-deploy --patch "$(sed "s/\${tag}/${{steps.meta.outputs.version}}/" < ./config/patch.yml)"`，国内可使用**Coding持续部署**或**云效流水线**等。

## License

[BSD-3-Clause License](https://github.com/weplanx/api/blob/main/LICENSE)