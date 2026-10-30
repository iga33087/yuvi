<template>
  <div id="app" class="bg-light bg-gradient">
    <Loading v-if="store.loading" />
    <Header />
    <div class="content">
      <router-view class="contentBox" />
    </div>
  </div>
</template>

<script setup>
import { useStore } from './store/index.js'
import Loading from './components/Loading.vue'
import Header from './components/Header.vue'

const store = useStore();
</script>
