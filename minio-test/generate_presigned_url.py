from minio import Minio
from datetime import timedelta

client = Minio(
    "minio:9000",
    access_key="myminioadmin",
    secret_key="password@123456",
    secure=False
)

url = client.presigned_put_object(
    "firmware",
    "testfile.txt",
    expires=timedelta(minutes=60)
)

print(url)
