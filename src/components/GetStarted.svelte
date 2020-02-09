<script>
import ProjectPreview from './ProjectPreview.svelte'
import Modal from './Modal.svelte'
import Button from './Button.svelte'
import CreateProjectPopup from './CreateProjectPopup.svelte'

let projects = []

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

  {#each projects as project}
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