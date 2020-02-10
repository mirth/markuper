
<script>
import Layout from './components/Layout.svelte';
import SideBar from './components/SideBar.svelte'

export let project;

</script>


<Layout>
  <div slot='sidebar'>
    <SideBar />
  </div>
  <div slot='content'>
    <h1>{project.name}</h1>
  </div>
</Layout>
