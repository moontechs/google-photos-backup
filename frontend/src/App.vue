<template>
  <div class="row layout-top-spacing">
      <div class="col-xl-4 col-lg-12 col-md-12 col-sm-12 col-12 layout-spacing">
          <AddNewClient></AddNewClient>
      </div>
      
  </div>
</template>

<script>
import AddNewClient from './components/AddNewClient.vue'

export default {
  name: 'App',
  components: {
    AddNewClient
  }
}
</script>

<!-- <style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style> -->
