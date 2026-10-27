<template>
  <div class="cardBox">
    <div class="cardBoxTitle fs-3" v-if="title">{{title}}</div>
    <div class="cardBoxContent">
      <slot name="content" />
    </div>
  </div>
</template>

<script setup>
const props = defineProps(['title'])

</script>