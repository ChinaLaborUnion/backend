SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go

docker login --username=嘎嘎哩satsun registry.cn-shenzhen.aliyuncs.com
880719.a
docker build -t registry.cn-shenzhen.aliyuncs.com/satsun/china_labor_union_backend:backend .
docker push registry.cn-shenzhen.aliyuncs.com/satsun/china_labor_union_backend:backend