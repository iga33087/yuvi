<template>
  <div class="result">
    <div class="row">
      <VideoEdit :src="videoUrl1" />
      <div class="col-8 mb-3">
        <video ref="mediaBox" class="resultVideoBox rounded-2" @loadedmetadata="onLoaded" @error="onError" />
      </div>
      <div class="col-4 mb-3">
          <CardBox title="影片資訊" class="rounded-2">
            <template v-slot:content>
              <div class="list-group mb-2">
                <div class="list-group-item">
                  <div class="fs-6 text-secondary mb-2">媒體名稱</div>
                  <div class="text-end">{{mediaInfo?.name}}</div>
                </div>
                <div class="list-group-item">
                  <div class="fs-6 text-secondary mb-2">媒體大小</div>
                  <div class="text-end">{{global.formatBytes(mediaInfo?.size)}}</div>
                </div>
                <div class="list-group-item">
                  <div class="fs-6 text-secondary mb-2">媒體類型</div>
                  <div class="text-end">{{mediaInfo?.mimeType}}</div>
                </div>
                <div class="list-group-item">
                  <div class="fs-6 text-secondary mb-2">修改時間</div>
                  <div class="text-end">{{global.timeFormat(mediaInfo?.modTime,"YYYY/MM/DD HH:mm")}}</div>
                </div>
              </div>
            </template>
          </CardBox>
      </div>
      <div class="col-12">
        <div class="resultControlBar">
          <div class="d-flex align-items-center justify-content-center">
            <button type="button" class="btn btn-outline-secondary rounded-pill me-2" @click="mediaBox.play()">播放</button>
            <button type="button" class="btn btn-outline-secondary rounded-pill me-2" @click="mediaBox.pause()">暫停</button>
            <button type="button" class="btn btn-outline-secondary rounded-pill me-2">範圍剪裁</button>
            <button type="button" class="btn btn-outline-secondary rounded-pill me-2">播放速度</button>
            <button type="button" class="btn btn-outline-secondary rounded-pill me-2">循環播放</button>
            <button type="button" class="btn btn-outline-secondary rounded-pill me-2" @click="exportMedia">匯出</button>
          </div>
          <el-slider v-model="timeRange" range :step="0.01" :max="mediaDuration" :marks="mediaMark" @input="timeRangeInput" />
        </div>
      </div>
    </div>
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

const mediaBox=ref(null)
const mediaInfo=ref(null)
const mediaURL=ref(null)
const timeRange=ref([0,0])
const lastTimeRange=ref([0,0])
const mediaDuration=ref(0)
const videoUrl1=ref("")
const mediaMark=computed(()=> {
  if(!mediaDuration.value) return {}
  let num=0
  let res={}
  while(num<=mediaDuration.value) {
    res[num] = num / 60 +"min"
    num+=60
  }
  return res
})

async function getInit() {
    if(!route?.query?.path) return
    mediaInfo.value=await GetFileInfo(route.query.path)
    const videoUrl =`/api/files?path=${encodeURIComponent(route.query.path)}`
    mediaBox.value.src = videoUrl
    mediaBox.value.load()
    videoUrl1.value=videoUrl
}

function timeRangeInput(e) {
  if(e[0]!=lastTimeRange.value[0]) mediaBox.value.currentTime=e[0]
  else if(e[1]!=lastTimeRange.value[1]) mediaBox.value.currentTime=e[1]
  lastTimeRange.value=e
}

async function exportMedia() {
  let res=await SaveFileDialog()
  console.log(1111,res)
}

function onLoaded() {
    console.log("影片載入成功")
    console.log("影片長度:", mediaBox)
    mediaDuration.value=mediaBox.value.duration
}

function onError() {
    const video = mediaBox.value

    console.error("影片載入失敗:", {
        code: video.error?.code,
        message: video.error?.message,
        networkState: video.networkState,
        readyState: video.readyState,
        currentSrc: video.currentSrc
    })
}

onMounted(async ()=> {
    await getInit()
})
</script>