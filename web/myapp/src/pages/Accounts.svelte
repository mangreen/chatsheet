<script>
    import { onMount } from 'svelte';
    // 引入 authService 和 getAuthToken
    import { authService, getAuthToken } from '../api'; 
    import { navigate } from 'svelte5-router';

    // --- 應用程式狀態 ---
    let loading = true;
    let error = '';
    let accounts = [];
    
    // --- 連線表單狀態 ---
    let username = '';
    let password = '';
    let accessToken = ''; // li_at cookie
    let userAgent = '';   // 瀏覽器 User-Agent
    let connectType = 'basic'; // 'basic' 或 'cookie'

    // --- Checkpoint 狀態 ---
    let isCheckpoint = false;
    let checkpointType = ''; // '2FA', 'OTP', 'IN_APP_VALIDATION', etc.
    let checkpointAccountId = ''; // 用於解決 Checkpoint 的 ID
    let checkpointCode = ''; // 用戶輸入的驗證碼或手機號碼

    // --- 輔助變數 ---
    let connectError = '';
    let connectLoading = false;

    onMount(async () => {
        // 檢查是否有 token，如果沒有，導向登入頁面
        if (!getAuthToken()) {
            navigate('/login');
            return;
        }
        await fetchAccounts();
    });

    /**
     * 獲取已連結的帳號列表
     */
    async function fetchAccounts() {
        loading = true;
        error = '';
        try {
            const response = await authService.getAccounts();
            accounts = response.data.accounts;
        } catch (e) {
            const status = e.response?.status;
            if (status === 401) {
                // 如果是 401，攔截器會清除 token，我們直接導航
                navigate('/login');
                return;
            }
            error = e.response?.data?.error || e.message;
        } finally {
            loading = false;
        }
    }
    
    /**
     * 重設所有連線相關的狀態
     */
    function resetConnectionState() {
        connectError = '';
        connectLoading = false;
        username = '';
        password = '';
        accessToken = '';
        userAgent = '';
        isCheckpoint = false;
        checkpointType = '';
        checkpointAccountId = '';
        checkpointCode = '';
    }

    /**
     * 處理基本登入 (Username/Password) 或 Cookie 登入
     */
    async function handleConnect() {
        connectLoading = true;
        connectError = '';

        try {
            let response;
            if (connectType === 'basic') {
                response = await authService.connectLinkedInBasic(username, password);
            } else if (connectType === 'cookie') {
                response = await authService.connectLinkedInCookie(accessToken, userAgent);
            }

            // 成功連接：200 OK
            alert(`成功連結! Account ID: ${response.data.account_id}`);
            resetConnectionState();
            await fetchAccounts(); // 重新載入帳號列表

        } catch (e) {
            connectLoading = false;
            
            // 處理 Checkpoint (202 Accepted) 響應
            if (e.response && e.response.status === 202) {
                const data = e.response.data;
                isCheckpoint = true;
                checkpointAccountId = data.account_id;
                checkpointType = data.checkpoint_type;
                connectError = `需要解決 Checkpoint: ${checkpointType}. 請在 5 分鐘內輸入驗證碼。`;
            } else {
                // 處理其他錯誤
                connectError = e.response?.data?.details || e.response?.data?.error || e.message || '連線失敗';
            }
        } finally {
            if (!isCheckpoint) {
                connectLoading = false;
            }
        }
    }

    /**
     * 處理 Checkpoint 解決方案
     */
    async function handleCheckpointSolve() {
        connectLoading = true;
        connectError = '';

        try {
            const response = await authService.solveCheckpoint(checkpointAccountId, checkpointCode);

            // 成功連接：200 OK
            alert(`Checkpoint 解決成功! Account ID: ${response.data.account_id}`);
            resetConnectionState();
            await fetchAccounts(); // 重新載入帳號列表

        } catch (e) {
            connectLoading = false;
            
            // 處理新的 Checkpoint (202 Accepted) 響應
            if (e.response && e.response.status === 202) {
                const data = e.response.data;
                // 更新 Checkpoint 狀態，繼續等待新輸入
                checkpointAccountId = data.account_id; 
                checkpointType = data.checkpoint_type;
                checkpointCode = ''; // 清空輸入，等待新碼
                connectError = `新的 Checkpoint: ${checkpointType}. 請重新輸入。`;
            } else {
                 // 處理超時 (408) 或其他錯誤
                resetConnectionState(); // 關閉 Checkpoint 介面
                alert('Checkpoint 解決失敗或已超時。請重新開始連線流程。');
                connectError = e.response?.data?.details || e.response?.data?.error || e.message || '解決失敗';
            }
        }
    }

    /**
     * 處理登出
     */
    async function handleLogout() {
        await authService.logout();
        navigate('/login');
    }
</script>

<div class="accounts-container">
    <header>
        <h1>Unipile Accounts</h1>
        <button on:click={handleLogout} class="btn-logout">Logout</button>
    </header>

    <hr>

    {#if loading}
        <p class="loading">Loading...</p>
    {:else if error}
        <p class="error-message">Error: {error}</p>
    {:else}
        <section class="connect-form-section">
            <h2>LinkedIn</h2>

            {#if isCheckpoint}
                <form on:submit|preventDefault={handleCheckpointSolve} class="connect-form checkpoint-form">
                    <p class="checkpoint-info">Checkpoint type: <strong>{checkpointType}</strong></p>
                    <label>
                        {checkpointType === '2FA' || checkpointType === 'OTP' ? 'Verify (Code)' : checkpointType === 'PHONE_REGISTER' ? 'Mobile (+contry code)' : 'Input'}
                        <input type="text" bind:value={checkpointCode} required disabled={connectLoading}>
                    </label>
                    
                    {#if connectError}
                        <p class="error-message">{connectError}</p>
                    {/if}

                    <button type="submit" disabled={connectLoading}>
                        {connectLoading ? 'Submitted...' : 'Solved Checkpoint'}
                    </button>
                    <button type="button" on:click={resetConnectionState} class="btn-cancel">Cancel</button>
                </form>

            {:else}
                <div class="connect-tabs">
                    <button 
                        class:active={connectType === 'basic'} 
                        on:click={() => { connectType = 'basic'; connectError = ''; }}
                    >
                        Basic
                    </button>
                    <button 
                        class:active={connectType === 'cookie'} 
                        on:click={() => { connectType = 'cookie'; connectError = ''; }}
                    >
                        Cookie
                    </button>
                </div>
                
                <form on:submit|preventDefault={handleConnect} class="connect-form">
                    {#if connectType === 'basic'}
                        <label>
                            <input type="email" placeholder="E-Mail" bind:value={username} required disabled={connectLoading}>
                        </label>
                        <label>
                            <input type="password" placeholder="Password" bind:value={password} required disabled={connectLoading}>
                        </label>
                    {:else}
                        <label>
                            <input type="text" placeholder="Access Token" bind:value={accessToken} required disabled={connectLoading}>
                        </label>
                    {/if}

                    {#if connectError}
                        <p class="error-message">{connectError}</p>
                    {/if}

                    <button type="submit" disabled={connectLoading}>
                        {connectLoading ? 'Connecting...' : 'Connect'}
                    </button>
                </form>
            {/if}

        </section>

        <hr>

        <section class="linked-accounts">
            <h2>Accounts: ({accounts.length})</h2>
            {#if accounts.length === 0}
                <p>No accounts is found.</p>
            {:else}
                <ul class="account-list">
                    {#each accounts as account (account.account_id)}
                        <li class="account-item">
                            <span class="provider">{account.provider.toUpperCase()}</span>
                            <span class="id-display">ID: {account.account_id}</span>
                            </li>
                    {/each}
                </ul>
            {/if}
        </section>
    {/if}
</div>

<style>
    .accounts-container {
        max-width: 800px;
        margin: 40px auto;
        padding: 20px;
        background-color: #f4f4f9;
        border-radius: 8px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px;
    }

    h1 {
        color: #333;
        font-size: 2em;
    }

    h2 {
        color: #555;
        border-bottom: 2px solid #ddd;
        padding-bottom: 10px;
        margin-bottom: 20px;
    }

    hr {
        border: none;
        border-top: 1px dashed #ccc;
        margin: 30px 0;
    }

    .btn-logout {
        background-color: #e74c3c;
        color: white;
    }

    /* --- List Style --- */
    .account-list {
        list-style: none;
        padding: 0;
    }

    .account-item {
        display: flex;
        justify-content: space-between;
        padding: 10px 15px;
        margin-bottom: 8px;
        background-color: white;
        border: 1px solid #eee;
        border-radius: 4px;
        align-items: center;
    }
    
    .provider {
        font-weight: bold;
        color: #0077b5; /* LinkedIn Blue */
    }

    .id-display {
        font-size: 0.9em;
        color: #777;
    }

    /* --- Form Style --- */
    .connect-form, .checkpoint-form {
        display: flex;
        flex-direction: column;
        gap: 15px;
        padding: 20px;
        border: 1px solid #ddd;
        border-radius: 4px;
        background-color: #fff;
    }
    
    label {
        display: flex;
        flex-direction: column;
        font-weight: 500;
        color: #555;
    }

    input[type="text"], input[type="email"], input[type="password"] {
        padding: 10px;
        margin-top: 5px;
        border: 1px solid #ccc;
        border-radius: 4px;
        font-size: 1em;
    }

    /* --- Buttons --- */
    button {
        padding: 10px 15px;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-weight: bold;
        transition: background-color 0.3s;
    }

    button[type="submit"] {
        background-color: #0077b5;
        color: white;
    }

    button[type="submit"]:hover {
        background-color: #005f91;
    }

    button:disabled {
        background-color: #ccc;
        cursor: not-allowed;
    }
    
    .btn-cancel {
        background-color: #95a5a6;
        color: white;
    }
    
    .btn-cancel:hover {
        background-color: #7f8c8d;
    }
</style>