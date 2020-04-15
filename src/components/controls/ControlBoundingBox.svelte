<script>
import Table from 'svelte-atoms/Table/Table.svelte';
import Tbody from 'svelte-atoms/Table/Tbody.svelte';
import { sampleView, sampleMarkup, assessState } from '../../store';
import ControlList from './ControlList.svelte';
import BoxMarkup from './BoxMarkup.svelte';

export let field;

let upperLeft;
let downRight;

$: if (!$sampleMarkup[field.group]) {
  $sampleMarkup[field.group] = [];
}

$: img = $assessState.imageElement;

function cornersToBox() {
  let width = downRight.x - upperLeft.x;
  let height = downRight.y - upperLeft.y;

  let x = Math.min(Math.max(0, upperLeft.x), img.width);
  let y = Math.min(Math.max(0, upperLeft.y), img.height);

  const scale = img.naturalWidth / img.clientWidth;
  x *= scale;
  y *= scale;
  width *= scale;
  height *= scale;

  return {
    x,
    y,
    width,
    height,
  };
}

function computePos(ev) {
  const rect = img.getBoundingClientRect();

  return {
    x: ev.pageX - rect.left + img.offsetLeft,
    y: ev.pageY - rect.top + img.offsetTop,
  };
}

function handleMousedown(ev) {
  if ($assessState.focusedGroup !== field.group) {
    return;
  }

  if (ev.x < img.x || ev.x > img.x + img.width) {
    return;
  }

  if (ev.y < img.y || ev.y > img.y + img.height) {
    return;
  }

  upperLeft = computePos(ev);
}

function handleMouseup() {
  if (!(upperLeft && downRight)) {
    return;
  }
  if (($assessState.markup.box.width === 0) || ($assessState.markup.box.height === 0)) {
    return;
  }

  [$assessState.focusedGroup] = field.fields_order;

  upperLeft = null;
  downRight = null;
}

$: if (upperLeft && downRight) {
  $assessState.markup.box = cornersToBox();
}

function handleMousemove(ev) {
  if (upperLeft) {
    downRight = computePos(ev);

    if (ev.x > img.x + img.width) {
      downRight.x = img.width;
    }
    if (ev.y > img.y + img.height) {
      downRight.y = img.height;
    }
  }

  if (upperLeft && downRight) {
    const box = cornersToBox();
    $assessState.activeBBox = box;
  }
}

function removeBox(boxIdx) {
  return () => {
    $sampleMarkup[field.group] = $sampleMarkup[field.group].filter((_el, i) => i !== boxIdx);
  };
}

function selectBox(i) {
  $sampleView.selectedBox = i;
}

</script>

<svelte:window
  on:mousedown={handleMousedown}
  on:mouseup={handleMouseup}
  on:mousemove={handleMousemove}
  />


<ControlList owner={field} />

<div id='boxes'>
  <Table>
    <Tbody>
      {#each $sampleMarkup[field.group] as markup, i}
        <tr>
          <td>
            <BoxMarkup
              {markup}
              isSelected={$sampleView.selectedBox === i}
              on:mouseover={() => selectBox(i)}
              onCrossPressed={removeBox(i)}
              />
          </td>
        </tr>
      {/each}
    </Tbody>
  </Table>
</div>

<style>

</style>