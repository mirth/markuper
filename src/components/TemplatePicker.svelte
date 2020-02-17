<script>
import { onMount } from 'svelte';
import api from '../api';
import _ from 'lodash';
import TemplatePreview from './TemplatePreview.svelte';

export let selectedTemplate;

let selectedTemplateTask = null;
let templateList = [];

onMount(async () => {
  const res = await api.get('/project_templates');
  templateList = res.templates;
})



$: selectedTemplate.template = Object.assign({}, _.find(
  templateList,
  {'task': selectedTemplateTask}),
)
</script>

<style>

.grid {
  display: grid;
}

</style>


<div class="grid">
  {#each templateList as template}
    <label>
      <input type=radio bind:group={selectedTemplateTask} value={template.task} />
    </label>
  {/each}
</div>

{#if selectedTemplate.template.task}
  <TemplatePreview template={selectedTemplate} />
{/if}
