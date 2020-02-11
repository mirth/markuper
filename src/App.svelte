<script>
import './global.svelte';

import PageBlank from './components/PageBlank.svelte';
import GetStarted from './components/GetStarted.svelte';
import Project from './components/Project.svelte';
import Router from 'svelte-spa-router'

const routes = {
    '/': GetStarted,
    '/project/:project_id': Project,
}
</script>

<Router {routes}/>