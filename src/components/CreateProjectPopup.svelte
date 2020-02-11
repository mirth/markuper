<script>
import { newProject, projects, fetchProjectList } from "../store.js";
import api from '../api';
import ProjectTemplatePicker from './ProjectTemplatePicker.svelte'

export let close;

let name = '';

async function createNewProject() {
  await api.post('/project', $newProject);
  await fetchProjectList();
  close();
}

</script>

<form on:submit|preventDefault={createNewProject}>
  <input bind:value={$newProject.description.name} placeholder="New project" minlength="1">
  <ProjectTemplatePicker />
  <button>Create</button>
</form>