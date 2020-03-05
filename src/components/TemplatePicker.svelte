<script>
import { onMount } from 'svelte';
import _ from 'lodash';
import Select from 'svelte-atoms/Select.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import api from '../api';
import TemplatePreview from './TemplatePreview.svelte';

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
};

$: options = _.map(templateList, (t) => ({
  label: t.task, value: t.task,
}));

</script>

<Select bind:value={selectedTemplateTask} {options} title='Select project Task' placeholder='Select task' />
<Spacer size={16} />
{#if selectedTemplate.template.task}
  <TemplatePreview template={selectedTemplate} />
{/if}
