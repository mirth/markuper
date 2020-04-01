<script>
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Block from 'svelte-atoms/Block.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Popup from 'svelte-atoms/Popup.svelte';
import Button from 'svelte-atoms/Button.svelte';
import Container from 'svelte-atoms/Grids/Container.svelte';
import CreateProjectPopup from './CreateProjectPopup.svelte';
import PageBlank from './PageBlank.svelte';
import ProjectList from './ProjectList.svelte';


let isNewProjectPopupShown = false;

function showNewProjectPopup() {
  isNewProjectPopupShown = true;
}

function closeNewProjectPopup() {
  isNewProjectPopupShown = false;
}

</script>

<Container>
<Row>
<Cell xs={4}><Block></Block></Cell>
<Cell xs={4}><Block>
  <Spacer size={128} />
  <div>
    <Button
      on:click={showNewProjectPopup}
      iconRight='chevron-right'
      style='display: inline;'
    >
      Create new project
    </Button>
  </div>
</Block></Cell>
<Cell xs={4}><Block></Block></Cell>
</Row>
<ProjectList />
</Container>

<Popup isOpen={isNewProjectPopupShown} on:close={closeNewProjectPopup}>
  <div style='min-width: 600px;'>
    <CreateProjectPopup close={closeNewProjectPopup} />
  </div>
</Popup>

<style>
div {
  text-align: center;
}
</style>