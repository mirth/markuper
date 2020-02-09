<script>
import api from '../api';
import { newProject } from "../store.js";

let name = '';

async function createNewProject() {
  const res = await api.post('/project', $newProject);
}

</script>

<form on:submit|preventDefault={createNewProject}>
  <input bind:value={$newProject.name} placeholder="New project" minlength="1">
  <button>Create</button>
</form>