import axios from 'axios';

const API_BASE_URL = '/';

// ----------------------------------------------------
// 處理 JWT 存儲與傳輸
// ----------------------------------------------------

// 使用 localStorage 存儲 token
const TOKEN_KEY = 'authToken';

export function setAuthToken(token) {
    if (token) {
        localStorage.setItem(TOKEN_KEY, token);
    } else {
        localStorage.removeItem(TOKEN_KEY);
    }
}

export function getAuthToken() {
    return localStorage.getItem(TOKEN_KEY);
}

// ----------------------------------------------------
// 配置 Axios 實例
// ----------------------------------------------------

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// 設置請求攔截器，在每次發送請求前加入 Authorization Header
api.interceptors.request.use(config => {
    const token = getAuthToken();
    if (token) {
        // [NEW] 將 token 加入 Authorization: Bearer Header
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, error => {
    return Promise.reject(error);
});

// 設置響應攔截器，處理 401 Unauthorized 錯誤
api.interceptors.response.use(response => {
    return response;
}, error => {
    // 檢查是否為 401 錯誤，如果是，則清除 token 並導向登入頁面
    if (error.response && error.response.status === 401) {
        setAuthToken(null);
        // [TODO] 實際應用中，這裡應該觸發導航到登入頁面的邏輯
        console.error("401 Unauthorized, token cleared.");
    }
    return Promise.reject(error);
});

// ----------------------------------------------------
// 服務定義
// ----------------------------------------------------

export const authService = {
    login: async (email, password) => {
        const response = await api.post('/auth/login', { email, password });
        if (response.data && response.data.token) {
            setAuthToken(response.data.token);
        }
        return response.data;
    },

    signup: (password, email) => api.post('/auth/signup', { email, password }),
    
    logout: () => { 
        setAuthToken(null); 
    },

    getAccounts: () => api.get('/api/unipile'),
    
    connectLinkedInBasic: (username, password) => api.post('/api/unipile/linkedin/basic', { username, password }),
    
    connectLinkedInCookie: (accessToken, userAgent) => api.post('/api/unipile/linkedin/cookie', { access_token: accessToken, user_agent: userAgent }),
    
    solveCheckpoint: (accountId, code) => api.post('/api/unipile/linkedin/checkpoint', { account_id: accountId, code: code }),
};