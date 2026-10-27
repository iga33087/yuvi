<template>
  <div class="edit">
    <VideoEdit class="editVideo" :src="mediaURL" :info="mediaInfo" />
  </div>
</template>

<script setup>
import {ref,onMounted,nextTick,computed} from "vue"
import global from "../assets/js/global.js"
import CardBox from '../components/CardBox.vue'
import VideoEdit from '../components/VideoEdit.vue'
import { useRoute } from 'vue-router'
import { GetFileInfo,DownloadFileByPath,SaveFileDialog } from "../../wailsjs/go/main/App";
const route = useRoute()

const mediaInfo=ref(null)
const mediaURL=ref(null)

async function getInit() {
    if(!route?.query?.path) return
    mediaInfo.value=await GetFileInfo(route.query.path)
    const videoUr =`/api/files?path=${encodeURIComponent(route.query.path)}`
    mediaURL.value=videoUr
}

onMounted(async ()=> {
    await getInit()
})
</script>