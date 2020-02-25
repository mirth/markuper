<script>
import Popup from 'svelte-atoms/Popup.svelte';
import Button from 'svelte-atoms/Button.svelte';
import CreateProjectPopup from './CreateProjectPopup.svelte';
import ProjectPreview from './ProjectPreview.svelte';
import { projects } from '../store';
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

<Popup isOpen={isNewProjectPopupShown} on:close={closeNewProjectPopup}>
  <CreateProjectPopup close={closeNewProjectPopup} />

  <div slot="footer">
    <Button on:click={closeNewProjectPopup}>Ok</Button>
  </div>
</Popup>
</PageBlank>