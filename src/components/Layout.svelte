
<style>
#grid {
  display: grid;
  height: 100px;
  grid-template-columns: repeat(6, 1fr);
  grid-template-rows: 100px;
}
</style>

<div id="grid">
  <slot name="sidebar"></slot>
  <slot name="content"></slot>
</div>