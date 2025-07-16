# filemgr

```
go mod init XXX
go mod tidy
go build -o filemgr_backend cmd/main.go
./filemgr_backend
```

# How to test
```
# generate upload presigned url
curl -X POST http://localhost:8080/fileUpload/uploadUrl \
  -H "Content-Type: application/json" \
  -d '{
    "bucketName": "firmware",
    "fileName": "folder1/testfile.txt",
    "tokenDurationSeconds": 3600
  }'

{"presignedUrl":"http://minio:9000/firmware/testfile.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256..."}

# upload file via presigned url
curl -X PUT -T testfile.txt "$presignedUrl"
```
