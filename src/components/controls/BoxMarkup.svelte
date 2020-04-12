<script>
import IconButton from 'svelte-atoms/IconButton.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Block from "svelte-atoms/Block.svelte";
import Container from 'svelte-atoms/Grids/Container.svelte';
import { sampleMarkup } from '../../store';

export let markup;
export let isSelected;
export let onCrossPressed;


</script>

<Container>
  <Row>
    <Cell xs={11}>
      <div on:mouseover>
        <Block
          type={isSelected ? 'selected' : 'block1'}
          style='padding: 3px;'
        >
          <span><small>left: {Math.round(markup.box.x)}, top: {Math.round(markup.box.y)}, width: {markup.box.width}, height: {markup.box.height}</small></span>
        </Block>
      </div>
    </Cell>
    <Cell xs={1}>
      <IconButton
        on:click={onCrossPressed}
        icon='close'
        type='flat'
        status='primary'
        size='medium'
      />
    </Cell>
  </Row>
</Container>

<style>

/* div {
  background-color: var(--palette-primary-4);
} */

</style>