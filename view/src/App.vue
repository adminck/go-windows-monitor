<template>
  <div id="app">
    <barmonitor  />
  </div>
</template>

<script>
import barmonitor from './components/barmonitor.vue'

export default {
  name: 'App',
  components: {
    barmonitor
  }
}
</script>

<style>
#app {
  height: -webkit-fill-available;
}

body{
  margin: 0px;
  padding: 0px;
}

::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-thumb {
  background-color: #a1a3a9;
  border-radius: 8px;
}
</style>
