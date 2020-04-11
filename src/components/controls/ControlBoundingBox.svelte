<script>
import {sampleView, sampleMarkup, assessState} from '../../store';
import ControlList from './ControlList.svelte';
import { getFieldsInOrderFor } from '../../project';

export let field;

let upperLeft;
let downRight;
$: if(!$sampleMarkup[field.group]) {
  $sampleMarkup[field.group] = [];
}

$: img = $assessState.imageElement;

function cornersToBox(upperLeft, downRight) {
    const width = downRight.x - upperLeft.x;
    const height = downRight.y - upperLeft.y;

    return {
      x: upperLeft.x,
      y: upperLeft.y,
      width: width,
      height: height,
    }
}

function computePos(ev) {
  const rect = img.getBoundingClientRect()

  return {
    x: ev.pageX - rect.left + img.offsetLeft,
    y: ev.pageY - rect.top + img.offsetTop,
  }
}

function handleMousedown(ev) {
  if($assessState.focusedGroup !== field.group) {
    return;
  }

  const img = $assessState.imageElement;

  if(ev.x < img.x || ev.x > img.x + img.width) {
    return
  }

  if(ev.y < img.y || ev.y > img.y + img.height) {
    return
  }

  upperLeft = computePos(ev)
}

function handleMouseup(ev) {
  if(!(upperLeft && downRight)) {
    return;
  }
  if(($assessState.markup.box.width === 0) || ($assessState.markup.box.height === 0)) {
    return;
  }

  $assessState.focusedGroup = field.fields_order[0];

  upperLeft = null;
  downRight = null;
}

$: if(upperLeft && downRight) {
  $assessState.markup.box = cornersToBox(upperLeft, downRight);
}

function handleMousemove(ev) {
  if(upperLeft) {
    downRight = computePos(ev);

    const rect = img.getBoundingClientRect()

    if(ev.x > img.x + img.width) {
      downRight.x = img.width;
    }
    if(ev.y > img.y + img.height) {
      downRight.y = img.height;
    }
  }

  if(upperLeft && downRight) {
    const box = cornersToBox(upperLeft, downRight)
    $assessState.activeBBox = box;
  }
}

</script>

<svelte:window
  on:mousedown={handleMousedown}
  on:mouseup={handleMouseup}
  on:mousemove={handleMousemove}
  />


<ControlList owner={field} />
<ul id='boxes'>
  {#each $sampleMarkup[field.group] as mark, i}
  <li>
    left: {mark.box.x}, top: {mark.box.y}, width: {mark.box.width}, height: {mark.box.height}
  </li>
  {/each}
</ul>
