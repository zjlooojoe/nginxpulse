import axios from 'axios';

const ACCESS_KEY_STORAGE = 'nginxpulse_access_key';
const ACCESS_KEY_HEADER = 'X-NginxPulse-Key';
const ACCESS_KEY_EVENT = 'nginxpulse:access-key-required';

const client = axios.create({
  baseURL: '/',
  timeout: 15000,
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
  },
});

client.interceptors.request.use((config) => {
  const accessKey = localStorage.getItem(ACCESS_KEY_STORAGE);
  if (accessKey) {
    config.headers[ACCESS_KEY_HEADER] = accessKey;
  }
  return config;
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status;
    const message = error?.response?.data?.error || error?.message || '请求失败';
    if (status === 401) {
      window.dispatchEvent(
        new CustomEvent(ACCESS_KEY_EVENT, {
          detail: { message },
        })
      );
    }
    return Promise.reject(new Error(message));
  }
);

export default client;
