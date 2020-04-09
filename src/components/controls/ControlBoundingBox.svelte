<script>
import {sampleView, activeMarkup, assessState} from '../../store';
import ControlList from './ControlList.svelte';
import { getFieldsInOrderFor } from '../../project';

export let field;
export let prefix;

// const [fields, groupsOrder] = getFieldsInOrderFor(field);

let upperLeft;
let downRight;

let boxes = [];
$: $activeMarkup[field.group] = boxes
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
  if(ev.pageX) {

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
  const box = cornersToBox(upperLeft, downRight);
  if((box.width === 0) || (box.height === 0)) {
    return;
  }
  boxes = [...boxes, box];

  upperLeft = null;
  downRight = null;
  $assessState.activeBBox = null;
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

function handleKeydown(ev) {

}

function handleKeyup(ev) {

}

</script>

<svelte:window
  on:mousedown={handleMousedown}
  on:mouseup={handleMouseup}
  on:mousemove={handleMousemove}
  />

  <!-- on:keydown={handleKeydown}
  on:keyup={handleKeyup} -->

<ControlList templateLike={field} {prefix} />
