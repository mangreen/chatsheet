<script>
    import { authService } from '../api';
    import { navigate } from 'svelte5-router';

    let email = '';
    let password = '';
    let isRegister = false;
    let error = '';

    async function handleSubmit() {
        error = '';
        try {
            if (isRegister) {
                await authService.register(email, password);
                alert('註冊成功，請登入！');
                isRegister = false;
            } else {
                await authService.login(email, password);
                // 登入成功後導向帳號頁面
                navigate('/accounts');
            }
        } catch (e) {
            error = e.response?.data?.error || e.message; 
        }
    }
</script>

<div class="container">
    <h2>{isRegister ? '註冊' : '登入'}</h2>
    {#if error}
        <p style="color: red;">{error}</p>
    {/if}
    <form on:submit|preventDefault={handleSubmit}>
        <input type="email" bind:value={email} placeholder="Email" required />
        <input type="password" bind:value={password} placeholder="密碼" required />
        
        <button type="submit">{isRegister ? '註冊' : '登入'}</button>
    </form>
    
    <button on:click={() => isRegister = !isRegister} class="toggle-btn">
        {isRegister ? '已有帳號？去登入' : '還沒有帳號？去註冊'}
    </button>
</div>

<style>
    .container { max-width: 400px; margin: 50px auto; padding: 20px; border: 1px solid #ccc; }
    input { width: 100%; padding: 10px; margin-bottom: 10px; box-sizing: border-box; }
    button { width: 100%; padding: 10px; background-color: #007bff; color: white; border: none; cursor: pointer; }
    .toggle-btn { margin-top: 10px; background-color: #6c757d; }
</style>