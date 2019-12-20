<style type="text/scss">
.small {
  --radius: .75em;
}

.protrude {
  background-color: #E4ECF6;
  box-shadow: 5px 5px 10px var(--shadow-color), 1px 1px 2px var(--shadow-color), -6px -6px 15px rgba(255,255,255,0.9), -1px -1px 3px 1px rgba(255,255,255,0.9);
}
.protrude {
  border-radius: var(--radius);
}

.button {
  border: none;
  display: inline-block;
  padding: calc(var(--padding-small) - .15em);
  margin-left: var(--padding-small);

  margin-bottom: 40px;
  font-weight: bold;
  font-family: inherit;
}

</style>


<button on:click class="small protrude button">
  <slot></slot>
</button>