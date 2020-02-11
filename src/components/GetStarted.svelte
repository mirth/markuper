<script>
import api from '../api';
import ProjectPreview from './ProjectPreview.svelte'
import Modal from './Modal.svelte'
import Button from './Button.svelte'
import CreateProjectPopup from './CreateProjectPopup.svelte'
import { projects } from "../store.js";
import PageBlank from './PageBlank.svelte';

let isNewProjectPopupShown = false;

function showNewProjectPopup() {
  isNewProjectPopupShown = true;
}

function closeNewProjectPopup() {
  isNewProjectPopupShown = false;
}

</script>

<PageBlank>
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
    <p><CreateProjectPopup close={closeNewProjectPopup}/></p>
  </Modal>
{/if}
</PageBlank>