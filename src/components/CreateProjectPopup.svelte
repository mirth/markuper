<script>
import { fetchProjectList } from '../store';
import api from '../api';
import TemplatePicker from './TemplatePicker.svelte';
import DataSource from './DataSource.svelte';

export let close;

let projectName = '';
const selectedTemplate = {
  template: null,
};
const dataSource = {
  type: 'local_directory',
  source_uri: '',
};

async function createNewProject() {
  await api.post('/project', {
    description: {
      name: projectName,
    },
    template: selectedTemplate.template,
    data_source: dataSource,
  });
  await fetchProjectList();
  close();
}

</script>

<form>
  <input bind:value={projectName} placeholder="New project" minlength="1">
  <TemplatePicker {selectedTemplate} />
  <DataSource {dataSource} />
  <button on:click={createNewProject}>Create</button>
</form>
