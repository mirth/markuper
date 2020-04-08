<script>
import { onMount } from 'svelte';
import {sampleView, activeMarkup, assessState} from '../../store';

export let field;

let upperLeft;
let downRight;

onMount(() => {

});
let boxes = [];
$: $activeMarkup[field.group] = boxes


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
  const rect = $assessState.imageElement.getBoundingClientRect()

  return {
    x: ev.pageX - rect.left + $assessState.imageElement.offsetLeft,
    y: ev.pageY - rect.top + $assessState.imageElement.offsetTop,
  }
}

function handleMousedown(ev) {
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
  // if(!(upperLeft && downRight)) {
  //   return;
  // }
  // console.log('ev: ', ev)
  if(upperLeft) {
    downRight = computePos(ev);
  }
  // console.log('ev: ', ev)

  // const ctx = $sampleView.getContext('2d');
  // ctx.clearRect(0, 0, 320, 320)
  // console.log('upperLeft: ', upperLeft)
  // console.log('downRight: ', downRight)
  if(upperLeft && downRight) {
    const box = cornersToBox(upperLeft, downRight)
    $assessState.activeBBox = box;
  //   drawBox(ctx, box);
  }

  // for(let box of boxes) {
  //   drawBox(ctx, box);
  // }
  // ctx.clear();
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
  on:keydown={handleKeydown}
  on:keyup={handleKeyup}
  />

KEK
<!-- <canvas
  bind:this={canvas}
></canvas> -->