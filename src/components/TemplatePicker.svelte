<script>
import { onMount } from 'svelte';
import _ from 'lodash';
import Select from 'svelte-atoms/Select.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Input from 'svelte-atoms/Input.svelte';
import api from '../api';
import TemplatePreview from './TemplatePreview.svelte';

export let selectedTemplate;

let selectedTemplateTask = 'classification';
let templateList = [];

onMount(async () => {
  const res = await api.get('/project_templates');
  templateList = res.templates;
  onTemplateSelected({detail: selectedTemplateTask})
});

function onTemplateSelected(ev) {
  selectedTemplateTask = ev.detail
  const template = {
    ..._.find(
      templateList,
      { task: selectedTemplateTask },
    ),
  }
  selectedTemplate.template = template
}

function isXMLValid(xml) {
  const oParser = new DOMParser();
  const oDOM = oParser.parseFromString(xml, 'text/xml');
  return oDOM.getElementsByTagName('parsererror').length === 0;
}

function validateXML(ev) {
  if(!isXMLValid(ev.target.value)) {
    selectedTemplate.error = "This doesn't look like valid xml"
  } else [
    selectedTemplate.error = null
  ]
}

$: options = _.map(templateList, (t) => ({
  label: t.task, value: t.task,
}));

</script>

<Input
  bind:value={selectedTemplate.template.xml}
  on:input={validateXML}
  textarea
  rows={5}
  />

{#if selectedTemplate.error}
<span>{selectedTemplate.error}</span>
{/if}

<Select
  on:change={onTemplateSelected}
  value={selectedTemplateTask}
  {options}
  title='Select project Task'
  placeholder='Select task'
  />

<style>
span {
  color: red;
}
</style>