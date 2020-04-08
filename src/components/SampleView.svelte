<script>
import {sampleView, activeMarkup, assessState} from '../store'
export let sample;

const field = sample.project.template.bounding_boxes[0];
let boxes = []
$: boxes = $activeMarkup[field.group] || []

// $: console.log('boxes: ', boxes)
// console.log('boxes: ', boxes)
</script>


<div class='image-container'>
  <img
    src='file://{sample.sample.image_uri}'
    alt='KEK'
    draggable=false
    bind:this={$assessState.imageElement}
    />

  {#each boxes as box}
    <div style={`
      width: ${box.width}px;
      height: ${box.height}px;
      left: ${box.x}px;
      top: ${box.y}px;
    `} class='box' />
  {/each}

  {#if $assessState.activeBBox}
    <div style={`
      width: ${$assessState.activeBBox.width}px;
      height: ${$assessState.activeBBox.height}px;
      left: ${$assessState.activeBBox.x}px;
      top: ${$assessState.activeBBox.y}px;
    `} class='box' />
  {/if}
  <!-- <canvas bind:this={$sampleView} >

  </canvas> -->
</div>

<style>
canvas {
  position:absolute;
  width: 100%;
  height: 100%;

  z-index:20;
}

img {
  /* position:absolute; */
  z-index:1;

  max-width: 100%;
  border: 1px solid black;

  display:block;
  margin-left:auto;
  margin-right:auto;
}

.box {
  position: absolute;
  border: 1px solid green;
}

/* .image-container {
  position:absolute;
  display:inline-block;
  /* padding: 0 45px 45px 0; */
  /* margin: 0 auto; */
  /* width: 100%; */
/* } */


</style>