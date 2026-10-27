<template>
  <div class="videoEdit">

  <el-dialog v-model="showExportDialog" title="匯出媒體" width="50vw">
    <el-form :model="form" label-width="auto" style="max-width: 600px">
      <el-form-item label="畫面幀率">
        <el-input-number v-model="exportConfig.fps" :min="0" :max="60" />
        <div class="ms-2">0為預設幀率</div>
      </el-form-item>
      <el-form-item label="畫面比例">
        <el-select v-model="exportConfig.ratio" style="width: 240px">
          <el-option label="原始比例" :value="null" />
          <el-option label="120p" :value="120" />
          <el-option label="240p" :value="240" />
          <el-option label="360p" :value="360" />
          <el-option label="480p" :value="480" />
          <el-option label="720p" :value="720" />
          <el-option label="1080p" :value="1080" />
        </el-select>
        <div class="ms-2">空值為原始比例</div>
      </el-form-item>
      <el-form-item label="最佳化調色">
        <el-switch v-model="exportConfig.palette" />
        <div class="ms-2">開啟此選項，匯出的容量會較大</div>
      </el-form-item>
    </el-form>
    <div class="d-flex align-items-center justify-content-end">
      <button type="button" class="btn btn-primary rounded-pill" @click="exportFunc">匯出</button>
    </div>
  </el-dialog>

    <div class="videoEditBox">
      <div class="videoEditBoxVideo">
        <div class="videoEditBoxVideoBox" ref="stageRef">
          <video class="videoEditBoxVideoBoxSrc" ref="mediaBox" :src="src" @loadedmetadata="onLoaded" @error="onError" @play="checkTime" />
          <v-stage :config="stageConfig" class="videoEditBoxVideoCrop" v-show="controlConfig.crop">
            <v-layer>
              <v-rect ref="cropRef" :config="cropConfig" @dragend="onDragend" @transformend="onTransformend" />
              <v-transformer ref="trRef" :config="trConfig" />
            </v-layer>
          </v-stage>
        </div>
      </div>
      <div class="videoEditBoxControl">
        <el-slider class="mb-4" v-model="timeRange" range :step="0.1" :max="mediaDuration" :marks="mediaMark" @input="timeRangeInput" />
        <div class="d-flex align-items-center justify-content-between">
          <div>
            <button type="button" class="btn rounded-pill me-2" :class="controlConfig.play ? 'btn-secondary':'btn-outline-secondary'" @click="controlConfig.play=!controlConfig.play">{{controlConfig.play ? '暫停':'播放'}}</button>
            <button type="button" class="btn rounded-pill me-2" :class="controlConfig.crop ? 'btn-secondary':'btn-outline-secondary'" @click="controlConfig.crop=!controlConfig.crop">範圍剪裁</button>
            <button type="button" class="btn rounded-pill me-2" :class="controlConfig.loop ? 'btn-secondary':'btn-outline-secondary'" @click="controlConfig.loop=!controlConfig.loop">循環播放</button>
            <button type="button" class="btn btn-outline-secondary rounded-pill me-2">播放速度</button>
          </div>
          <div>
            <button type="button" class="btn btn-outline-secondary rounded-pill" @click="showExportDialog=!showExportDialog">匯出</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref,computed,watch,onMounted} from "vue"
import { SaveFileDialog,APIExecMedia,ExportMedia } from "../../wailsjs/go/main/App";
const props = defineProps(['src','info'])

const showExportDialog=ref(false)
const exportConfig=ref({
  fps:0,
  ratio:null,
  palette:false,
})
const controlConfig = ref({
    play: false,
    crop: false,
    speed: 1,
    loop: true
})
const stageConfig = ref({ width: 0, height: 0 })
const trConfig = ref({
  rotateMode: 'corners',
})
const cropConfig = ref({
  x: 0, y: 0, width: 0, height: 0,
  stroke: '#0088ff', strokeWidth: 2,
  dash: [6, 3], draggable: true,
})

const mediaBox = ref(null)
const stageRef = ref(null)
const cropRef = ref(null)
const trRef = ref(null)
const timeRange = ref([0,0])
const lastTimeRange=ref([0,0])
const mediaDuration = ref(0)

onMounted(()=> {
  watch(()=>props.src,()=> {
    mediaBox.value.src = props.src
    mediaBox.value.load()
  },{
    immediate:true
  })
  watch(()=>controlConfig.value.play,()=> {
    if(controlConfig.value.play) mediaBox.value.play()
    else mediaBox.value.pause()
  })
  window.addEventListener("resize", () => {
    canvasInit()
  })
})

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

function canvasInit() {
  console.log(17)
  stageConfig.value.width = stageRef.value.offsetWidth
  stageConfig.value.height = stageRef.value.offsetHeight
  cropConfig.value.width = stageConfig.value.width
  cropConfig.value.height = stageConfig.value.height
  //const parentElement = stageRef.value.parentElement
  var cropNode = cropRef.value.getNode();
  var trNode = trRef.value.getNode();
  trNode.nodes([cropNode]);
  trNode.getLayer().batchDraw();
}

function timeRangeInput(e) {
  if(e[0]!=lastTimeRange.value[0]) mediaBox.value.currentTime=e[0]
  else if(e[1]!=lastTimeRange.value[1]) mediaBox.value.currentTime=e[1]
  lastTimeRange.value=e
}

function onLoaded() {
    console.log("影片載入成功")
    console.log("影片長度:", mediaBox)
    mediaDuration.value=mediaBox.value.duration
    timeRange.value[1]=mediaDuration.value
    lastTimeRange.value[1]=mediaDuration.value
    canvasInit()
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

function onDragend(e) {
  console.log(111,e)
}

function onTransformend(e) {
  console.log(112,e)
}

async function exportFunc() {
  let data= {
    info:props.info,
    controlConfig:controlConfig.value,
    cropConfig:cropConfig.value,
    timeRange:timeRange.value,
    exportConfig:exportConfig.value
  }

  let res=await ExportMedia(data)
  console.log(1111,res,data)
}

function checkTime() {
  let currentTime=0
  if (!mediaBox.value.paused && !mediaBox.value.ended) {
    currentTime=mediaBox.value.currentTime
    if(controlConfig.value.play && controlConfig.value.loop && currentTime >= timeRange.value[1]) mediaBox.value.currentTime=timeRange.value[0]
    if(!controlConfig.value.loop && currentTime >= timeRange.value[1]) controlConfig.value.play=false
  }
  requestAnimationFrame(checkTime);
}

</script>