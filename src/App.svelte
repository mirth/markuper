<script>
import './global.svelte';

import Router from 'svelte-spa-router'
import PageBlank from './components/PageBlank.svelte';
import GetStarted from './components/GetStarted.svelte';
import Project from './components/Project.svelte';
import AssessSample from './components/AssessSample.svelte';

const routes = {
    '/': GetStarted,
    '/project/:project_id': Project,
    '/project/:project_id/assess_sample': AssessSample,
}
</script>

<Router {routes}/>