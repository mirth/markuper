<script>
import { fetchProjectList } from '../store';
import api from '../api';
import ProjectTemplatePicker from './ProjectTemplatePicker.svelte';
import DataSource from './DataSource.svelte';

export let close;

let projectName = '';
const selectedTemplate = {};
const dataSource = {
  type: 'local_directory',
  source_uri: '',
};

async function createNewProject() {
  await api.post('/project', {
    description: {
      name: projectName,
    },
    template: selectedTemplate,
    data_source: dataSource,
  });
  await fetchProjectList();
  close();
}

</script>

<form on:submit|preventDefault={createNewProject}>
  <input bind:value={projectName} placeholder="New project" minlength="1">
  <ProjectTemplatePicker {selectedTemplate} />
  <DataSource {dataSource} />
  <button>Create</button>
</form>