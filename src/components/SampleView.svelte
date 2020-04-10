<script>
import {sampleView, activeMarkup, assessState} from '../store'
export let sample;

const field = sample.project.template.bounding_boxes[0];
let boxes = []
$: if(field) {
  boxes = $activeMarkup[field.group] || []
}

// $: console.log('boxes: ', boxes)
// console.log('boxes: ', boxes)
</script>


<div id='image-container'>
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

</div>

<style>

img {
  z-index:1;

  width: 100%;
  display:block;
  box-sizing: content-box;
}

.box {
  position: absolute;
  border: 1px solid green;

  box-sizing: border-box;
}

#image-container {
  position:relative;
  border: 1px solid black;
}



</style>