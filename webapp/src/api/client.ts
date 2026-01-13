import axios from 'axios';

const client = axios.create({
  baseURL: '/',
  timeout: 15000,
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
  },
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error?.response?.data?.error || error?.message || '请求失败';
    return Promise.reject(new Error(message));
  }
);

export default client;
