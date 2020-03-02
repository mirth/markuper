<script>
import { onMount } from 'svelte';
import _ from 'lodash';
import api from '../api';
import TemplatePreview from './TemplatePreview.svelte';
import Select from 'svelte-atoms/Select.svelte';

export let selectedTemplate;

let selectedTemplateTask = 'classification';
let templateList = [];

onMount(async () => {
  const res = await api.get('/project_templates');
  templateList = res.templates;
});


$: selectedTemplate.template = {
  ..._.find(
    templateList,
    { task: selectedTemplateTask },
  ),
}

$: options = _.map(templateList, (t) => {
  return {
    label: t.task, value: t.task
  }
});

$: console.log('options: ', options)
  // const options = [
  //   { label: "option 1", value: 1 },
  //   { label: "option 2", value: 2 }
  // ];
  // let value = 1;
</script>

<style>

/* .grid {
  display: grid;
} */

</style>


<!-- <div class="grid"> -->
  <!-- {#each templateList as template}
    <label>
      <input type=radio bind:group={c} value={template.task} />
    </label>
  {/each} -->
<!-- </div> -->
<Select bind:value={selectedTemplateTask} {options} title="Select Project Task" placeholder="Select task" />
<!-- <Select bind:value={} {options} title="Select" placeholder="Select option" /> -->



<!-- {#if selectedTemplate.template.task}
  <TemplatePreview template={selectedTemplate} />
{/if} -->
