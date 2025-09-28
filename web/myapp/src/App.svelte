<script>
    import { Router, Link, Route } from "svelte5-router";
    import Login from './pages/Login.svelte';
    import Accounts from './pages/Accounts.svelte';

    export let url = "";
</script>

<Router {url}>
  <nav>
    <Link to="/login">Login</Link> | 
    <Link to="/accounts">Accounts</Link>
  </nav>
  <main>
    <Route path="/" component={Login} />
    <Route path="/login" component={Login} />
    <Route path="/accounts" component={Accounts} />
  </main>
</Router>

<style>
    nav { padding: 10px; background: #f4f4f4; text-align: center; }
    nav a { margin: 0 10px; text-decoration: none; color: #333; }
</style>
