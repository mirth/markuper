<script>
import { sampleView, sampleMarkup, assessState } from '../store';

export let sample;

const field = sample.project.template.bounding_boxes[0];
let boxes = [];
$: if (field) {
  boxes = ($sampleMarkup[field.group] && $sampleMarkup[field.group]) || [];
}

function formatMarkup(markup) {
  const keys = Object.keys(markup).filter((k) => k !== 'box');

  return keys.map((k) => `${k}: ${markup[k]}`).join('\n');
}

</script>


<div id='image-container'>
  <img
    src='file://{sample.sample.image_uri}'
    alt='KEK'
    draggable=false
    bind:this={$assessState.imageElement}
    />

  {#each boxes as box, i}
    <div style={`
      width: ${box.box.width}px;
      height: ${box.box.height}px;
      left: ${box.box.x}px;
      top: ${box.box.y}px;
    `} class='box' class:box-selected={$sampleView.selectedBox === i}>
      {#if $sampleView.selectedBox === i}
        <span><small>{formatMarkup(box)}</small></span>
      {/if}
    </div>
  {/each}

  {#if $assessState.markup.box}
    <div style={`
      width: ${$assessState.markup.box.width}px;
      height: ${$assessState.markup.box.height}px;
      left: ${$assessState.markup.box.x}px;
      top: ${$assessState.markup.box.y}px;
    `} class='box'/>
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

.box-selected {
  border: 3px solid green;
}

#image-container {
  position:relative;
  border: 1px solid black;
}

span {
  background-color: green;
}


</style>