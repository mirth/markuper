<script>
import { link, push } from 'svelte-spa-router';
import Button from 'svelte-atoms/Button.svelte';
import { projects, fetchProject } from '../store';

function goToProject(projectId) {
  return async () => {
    await fetchProject(projectId);
    push(`/project/${projectId}`);
  };
}

</script>


<ul>
  <li>
    <a href={'/'} use:link>Home</a>
  </li>
  {#each $projects as project}
  <li>
    <Button type='empty' on:click={goToProject(project.project_id)}>{project.description.name}</Button>
  </li>
  {/each}
</ul>

<style>
ul {
  width: 3%;
}
</style>