<script>
import './global.svelte';
import 'svelte-atoms/Variables.svelte';

import Router from 'svelte-spa-router';
import GetStarted from './components/GetStarted.svelte';
import Project from './components/Project.svelte';
import AssessSample from './components/AssessSample.svelte';
import CreateProjectPage from './components/CreateProjectPage.svelte';

const routes = {
  '/': GetStarted,
  '/create_project': CreateProjectPage,
  '/project/:project_id': Project,
  '/project/:project_id/assess_sample': AssessSample,
  '/project/:project_id/samples/:sample_id': AssessSample,
};
</script>

<Router {routes}/>