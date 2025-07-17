<template>
  <div class="upload-container">
    <h1>升级包上传</h1>
    <form @submit.prevent="handleSubmit">
      <!-- 版本号输入 -->
      <div class="input-group">
        <label for="version">版本号</label>
        <input
          type="text"
          id="version"
          v-model="version"
          placeholder="例如：1.0.0"
          required
        />
      </div>

      <!-- 文件上传 -->
      <div class="input-group">
        <label for="file">选择升级包</label>
        <input
          type="file"
          id="file"
          @change="handleFileChange"
          accept=".zip,.tar.gz"
          required
        />
      </div>

     <!-- 上传进度和状态 -->
      <div v-if="progress > 0" class="progress-bar">
        <div class="progress" :style="{ width: progress + '%' }"></div>
        <span>{{ progress }}%</span>
      </div>
      <div v-if="status" class="status">{{ status }}</div>

      <button type="submit" :disabled="isLoading">
        {{ isLoading ? '处理中...' : '提交' }}
      </button>

    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';

const version = ref('');
const file = ref(null);
const isLoading = ref(false);
const progress = ref(0);
const status = ref('');

const handleFileChange = (e) => {
  file.value = e.target.files[0];
};

const handleSubmit = async () => {
  if (!version.value || !file.value) {
    alert('请填写版本号并选择文件');
    return;
  }
  isLoading.value = true;
  status.value = '开始处理文件...';
  progress.value = 0;

  try {
    status.value = '计算文件校验和...';
    const sha256 = await calculateSHA256(file.value);
    progress.value = 25;


    status.value = '获取升级包上传链接...';
    const firmwareUrl = await getPresignedURL({
      bucketName: 'firmware',
      objectKey: `${version.value}/${file.value.name}`,
      expireSeconds: 3600
    });
    const checksumUrl = await getPresignedURL({
      bucketName: 'firmware',
      objectKey: `${version.value}/${file.value.name}.sha256`,
      expireSeconds: 3600
    });
    progress.value = 50;

    status.value = '上传文件中...';
    await Promise.all([
      uploadFileToMinIO(firmwareUrl, file.value),
      uploadFileToMinIO(checksumUrl, new Blob([sha256], { type: 'text/plain' }))
    ]);
    progress.value = 100;

    status.value = '上传成功';
  } catch (err) {
    status.value = '上传失败';
    alert(`上传失败: ${err.message}`);
  } finally {
    isLoading.value = false;
  }
};

const calculateSHA256 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = async () => {
      try {
        const buffer = reader.result;
        const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        resolve(hashHex);
      } catch (err) {
        reject(err);
      }
    };
    reader.onerror = reject;
    reader.readAsArrayBuffer(file);
  });
};

const getPresignedURL = async (params) => {
  const response = await axios.post(
    '/api/file/uploadUrl',
    params,
    {
      headers: {
        'Content-Type': 'application/json'
      }
    }
  );
  return response.data.url;
};

const uploadFileToMinIO = async (url, file) => {
  const response = await axios.put(url, file, {
    headers: {
      'Content-Type': file.type || 'application/octet-stream'
    },
    onUploadProgress: (progressEvent) => {
      if (progressEvent.total) {
        progress.value = Math.min(
          progress.value + progressEvent.loaded / progressEvent.total * 20,
          90
        );
      }
    }
  });
};
</script>

<style scoped>
.upload-container {
  max-width: 500px;
  margin: 2rem auto;
  padding: 2rem;
  border: 1px solid #eee;
  border-radius: 8px;
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
}

.input-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

input[type="text"],
input[type="file"] {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

button {
  width: 100%;
  padding: 0.75rem;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
}

button:hover {
  background-color: #3aa876;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.progress-bar {
  margin: 1rem 0;
  height: 20px;
  background-color: #f0f0f0;
  border-radius: 4px;
  position: relative;
}

.progress {
  height: 100%;
  background-color: #42b983;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-bar span {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  color: #333;
  font-size: 0.8rem;
}

.status {
  margin: 1rem 0;
  text-align: center;
  color: #666;
}
</style>
