<script>
import Input from 'svelte-atoms/Input.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Button from 'svelte-atoms/Button.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import api from '../api';
import { goToProject } from '../project';
import PageBlank from './PageBlank.svelte';
import TemplatePicker from './TemplatePicker.svelte';
import DataSourcePicker from './DataSourcePicker.svelte';


let projectName = '';
const selectedTemplate = {
  template: {
    xml: '',
  },
  error: null,
};
const dataSources = {
  dataSources: [],
};

function isProjectValid() {
  if (projectNameError) {
    return false;
  }

  if (selectedTemplate.error) {
    return false;
  }

  if (dataSources.dataSources.length === 0) {
    return false;
  }

  return true;
}

let createProjectError;

async function createNewProject() {
  if (!isProjectValid()) {
    return;
  }

  const res = await api.post('/project', {
    description: {
      name: projectName,
    },
    template: selectedTemplate.template,
    data_sources: dataSources.dataSources,
  });

  if (res.status === 400) {
    createProjectError = res.data.error;
  }

  if (!Object.hasOwnProperty.call(res, 'status')) {
    createProjectError = null;
    await goToProject(res.project_id)();
  }
}

$: projectNameError = (projectName.trim().length === 0) && 'Project name should not be empty';

</script>

<PageBlank>
  <Row>
    <Cell xsOffset={2} xs={8}>
      <Typography type='headline' block>New project</Typography>
      <Spacer size={16} />
      <Input
        bind:value={projectName}
        name='projectName'
        title='Project name'
        size='big'
        placeholder='My cool project'
        error={projectNameError} />

      {#if projectNameError}
        <span>{projectNameError}</span>
      {/if}

      <Spacer size={16} />
      <TemplatePicker {selectedTemplate} />
      <Spacer size={24} />
      <DataSourcePicker {dataSources} />
      <Spacer size={36} />
      <div style='display: inline;'>
        {#if createProjectError}
          <span id='create_project_error'>{createProjectError}</span>
        {/if}
        <Button on:click={createNewProject} style='float: right; margin-bottom: 20px; margin-right: 20px;'>Create</Button>
      </div>
    </Cell>
  </Row>
</PageBlank>

<style>
span {
  color: red;
}
</style>