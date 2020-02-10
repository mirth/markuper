<script>
import { onMount } from 'svelte';
import api from '../api';
import ProjectPreview from './ProjectPreview.svelte'
import Modal from './Modal.svelte'
import Button from './Button.svelte'
import CreateProjectPopup from './CreateProjectPopup.svelte'
import { projects, fetchProjectList } from "../store.js";


onMount(fetchProjectList);

let isNewProjectPopupShown = false;

function showNewProjectPopup() {
  isNewProjectPopupShown = true;
}

function closeNewProjectPopup() {
  isNewProjectPopupShown = false;
}

</script>


<ul>
  <li>
    <Button on:click={showNewProjectPopup}>Create new project</Button>
  </li>

  {#each $projects as project}
    <li>
      <ProjectPreview project={project} />
    </li>
  {/each}
</ul>

{#if isNewProjectPopupShown}
  <Modal on:click={closeNewProjectPopup}>
    <p><CreateProjectPopup /></p>
  </Modal>
{/if}