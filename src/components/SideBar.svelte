<script>
import { projects } from '../store.js'

</script>

<style>
ul {
  width: 3%;
}
</style>

<ul>
  {#each $projects as project}

  <li>
    <a href={`/project/${project.project_id}`}>
      {project.description.name}
    </a>
  </li>
  {/each}
</ul>