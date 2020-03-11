<script>
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Block from 'svelte-atoms/Block.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Popup from 'svelte-atoms/Popup.svelte';
import Button from 'svelte-atoms/Button.svelte';
import CreateProjectPopup from './CreateProjectPopup.svelte';
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
<Row>
<Cell xs={4}><Block></Block></Cell>
<Cell xs={4}><Block>
  <Spacer size={128} />
  <Button on:click={showNewProjectPopup}>Create new project</Button>
</Block></Cell>
<Cell xs={4}><Block></Block></Cell>
</Row>

<Popup isOpen={isNewProjectPopupShown} on:close={closeNewProjectPopup}>
  <div style='min-width: 600px;'>
    <CreateProjectPopup close={closeNewProjectPopup} />
  </div>
</Popup>
</PageBlank>