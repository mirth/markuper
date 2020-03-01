<script>
import { fetchProjectList } from '../store';
import api from '../api';
import TemplatePicker from './TemplatePicker.svelte';
import DataSourcePicker from './DataSourcePicker.svelte';

export let close;

let projectName = '';
const selectedTemplate = {
  template: null,
};
const dataSources = [{
  type: 'local_directory',
  source_uri: '',
}];

async function createNewProject() {
  await api.post('/project', {
    description: {
      name: projectName,
    },
    template: selectedTemplate.template,
    data_sources: [dataSource],
  });
  await fetchProjectList();
  close();
}

</script>


<input bind:value={projectName} placeholder="New project" minlength="1">
<TemplatePicker {selectedTemplate} />
<DataSourcePicker {dataSources} />
<button on:click={createNewProject}>Create</button>
