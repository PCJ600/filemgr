<template>
  <div class="upload-container">
    <h1>固件升级包上传</h1>
    <form @submit.prevent="handleSubmit">
      <!-- 版本号输入 -->
      <div class="input-group">
        <label for="version">版本号</label>
        <input
          type="text"
          id="version"
          v-model="version"
          placeholder="例如：v1.0.0"
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

      <button type="submit">提交</button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';

const version = ref('');
const file = ref(null);
const isLoading = ref(false);

const handleFileChange = (e) => {
  file.value = e.target.files[0];
};

const handleSubmit = async () => {
  if (!version.value || !file.value) {
    alert('请填写版本号并选择文件');
    return;
  }
  isLoading.value = true;

  try {
    const presignedURL = await getPresignedURL({
      bucketName: 'firmware',
      fileName: `${version.value}/${file.value.name}`,
      tokenDurationSeconds: 3600
    });
    
    console.log(`预签名URL获取成功:\n${presignedURL}`);

    await uploadFileToMinIO(presignedURL, file.value);

  } catch (err) {
    console.error('上传失败:', err);
    alert(`操作失败: ${err.message}`);
  } finally {
    isLoading.value = false;
  }

};

const getPresignedURL = async (params) => {
  const response = await axios.post(
    '/api/fileUpload/uploadUrl',
    params,
    {
      headers: {
        'Content-Type': 'application/json'
      }
    }
  );
  return response.data.presignedUrl;
};

const uploadFileToMinIO = async (url, file) => {
  const response = await axios.put(url, file, {
    headers: {
      'Content-Type': 'application/octet-stream'
    }
  });
  console.log('文件上传完成', response);
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
</style>
