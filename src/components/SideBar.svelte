<script>
import { projects } from '../store.js'
import { link } from 'svelte-spa-router'

</script>

<style>
ul {
  width: 3%;
}
</style>

<ul>
  <li>
    <a href={`/`} use:link>Home</a>
  </li>
  {#each $projects as project}
  <li>
    <a href={`/project/${project.project_id}`} use:link>
      {project.description.name}
    </a>
  </li>
  {/each}
</ul>