<script>
import api from '../api';
import TemplatePreview from './TemplatePreview.svelte';
import { newProject } from "../store.js";

const templateList = api.get('/project_templates')

</script>

<style>

.grid {
  display: grid;
}

</style>


{#await templateList then list}
<div class="grid">
  {#each list.templates as template}
    <label>
      <input type=radio bind:group={$newProject.template} value={template}>
      <TemplatePreview {template} />
    </label>
  {/each}
</div>
{/await}