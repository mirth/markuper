<script>
import Input from 'svelte-atoms/Input.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Button from 'svelte-atoms/Button.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import api from '../api';
import { goToProject, enrichWithColor } from '../project';
import PageBlank from './PageBlank.svelte';
import TemplatePicker from './TemplatePicker.svelte';
import DataSourcePicker from './DataSourcePicker.svelte';


let projectName = '';
let selectedTemplate = {
  template: {
    xml: '',
  },
  error: null,
};

let dataSources = [];
let createProjectError;

$: projectNameError = (projectName.trim().length === 0) && 'Project name should not be empty';
$: isProjectValid = !(
  projectNameError
  || selectedTemplate.error
  || dataSources.length === 0
);

$: isCreatingProject = false;

$: isCreateProjectButtonEnabled = isProjectValid && !isCreatingProject;

async function createNewProject() {
  if (!isCreateProjectButtonEnabled) {
    return;
  }
  isCreatingProject = true;

  selectedTemplate.template.xml = enrichWithColor(selectedTemplate.template.xml);

  const res = await api.post('/project', {
    description: {
      name: projectName,
    },
    template: selectedTemplate.template,
    data_sources: dataSources,
    shuffle_samples: false, // fixme
  });

  if (res.status === 400) {
    createProjectError = res.data.error;
  }

  if (!Object.hasOwnProperty.call(res, 'status')) {
    createProjectError = null;
    await goToProject(res.project_id)();
  }

  isCreatingProject = false;
}

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
      <TemplatePicker bind:selectedTemplate />
      <Spacer size={24} />
      <DataSourcePicker bind:dataSources />
      <Spacer size={36} />
      <div style='display: inline;'>
        {#if createProjectError}
          <span id='create_project_error'>{createProjectError}</span>
        {/if}
        <Button
          on:click={createNewProject} style='float: right;'
          disabled={!isCreateProjectButtonEnabled}
          isLoading={isCreatingProject}
        >
          Create
        </Button>
      </div>
    </Cell>
  </Row>
</PageBlank>

<style>
span {
  color: red;
}
</style>