<script>
import { newProject, projects, fetchProjectList } from "../store.js";
import api from '../api';

let name = '';

async function createNewProject() {
  await api.post('/project', $newProject);
  await fetchProjectList();
}

</script>

<form on:submit|preventDefault={createNewProject}>
  <input bind:value={$newProject.name} placeholder="New project" minlength="1">
  <button>Create</button>
</form>