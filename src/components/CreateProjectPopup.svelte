<script>
import Input from 'svelte-atoms/Input.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import { fetchProjectList } from '../store';
import api from '../api';
import TemplatePicker from './TemplatePicker.svelte';
import DataSourcePicker from './DataSourcePicker.svelte';

export let close;

let projectName = '';
const selectedTemplate = {
  template: null,
};
const dataSources = {
  dataSources: []
};

async function createNewProject() {
  await api.post('/project', {
    description: {
      name: projectName,
    },
    template: selectedTemplate.template,
    data_sources: dataSources.dataSources,
  });
  await fetchProjectList();
  close();
}

</script>


<Input bind:value={projectName} title="Project name" size="big" placeholder="New project" />
<Spacer size={32} />
<TemplatePicker {selectedTemplate} />
<Spacer size={24} />
<DataSourcePicker {dataSources} />
<button on:click={createNewProject}>Create</button>
