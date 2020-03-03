<script>
import Input from 'svelte-atoms/Input.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Button from 'svelte-atoms/Button.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
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
  dataSources: [],
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

<Typography type='headline' block>New project</Typography>
<Spacer size={16} />
<Input bind:value={projectName} title="Project name" size="big" placeholder="My cool project" />
<Spacer size={16} />
<TemplatePicker {selectedTemplate} />
<Spacer size={24} />
<DataSourcePicker {dataSources} />
<Spacer size={36} />
<Button on:click={createNewProject} style='float: right; margin-bottom: 20px; margin-right: 20px;'>Create</Button>
