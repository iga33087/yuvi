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
              <v-rect ref="cropRef" :config="cropConfig" @click="updateCropDataConfig" @dragend="updateCropDataConfig" @transformend="updateCropDataConfig" />
              <v-transformer ref="trRef" :config="trConfig" />
            </v-layer>
          </v-stage>
        </div>
      </div>
      <div class="videoEditBoxControl">
        <el-slider class="mb-4" v-model="timeRange" range :step="0.1" :max="mediaDuration" :marks="mediaMark" @input="timeRangeInput" />
        <div class="d-flex align-items-center justify-content-between">
          <div class="d-flex align-items-center">
            <button type="button" class="btn d-flex align-items-center rounded-pill me-2" :class="controlConfig.play ? 'btn-secondary':'btn-outline-secondary'" @click="controlConfig.play=!controlConfig.play"><PauseIcon class="me-2" v-if="controlConfig.play" /><PlayIcon class="me-2" v-else /> {{controlConfig.play ? '暫停':'播放'}}</button>
            <button type="button" class="btn d-flex align-items-center rounded-pill me-2" :class="controlConfig.crop ? 'btn-secondary':'btn-outline-secondary'" @click="controlConfig.crop=!controlConfig.crop;testC()"><CropIcon class="me-2" />範圍剪裁</button>
            <button type="button" class="btn d-flex align-items-center rounded-pill me-2" :class="controlConfig.loop ? 'btn-secondary':'btn-outline-secondary'" @click="controlConfig.loop=!controlConfig.loop"><RepeatIcon class="me-2" />循環播放</button>
            <button type="button" class="btn d-flex align-items-center btn-outline-secondary rounded-pill me-2"><SpeedometerIcon class="me-2" />播放速度</button>
          </div>
          <div>
            <button type="button" class="btn d-flex align-items-center btn-outline-secondary rounded-pill" @click="showExportDialog=!showExportDialog"><FloppyIcon class="me-2" />匯出</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import PlayIcon from "../assets/icons/play-fill.svg"
import PauseIcon from "../assets/icons/pause-fill.svg"
import CropIcon from "../assets/icons/crop.svg"
import RepeatIcon from "../assets/icons/repeat.svg"
import SpeedometerIcon from "../assets/icons/speedometer2.svg"
import FloppyIcon from "../assets/icons/floppy.svg"
import {ref,computed,watch,onMounted,nextTick} from "vue"
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
const cropDataConfig = ref({
  x:0,
  y:0,
  width:0,
  height:0
})
const videoSizeConfig = ref({
  videoWidth:0,
  videoHeight:0,
  offsetWidth:0,
  offsetHeight:0,
  scaleX:1,
  scaleY:1
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
  stageConfig.value.width = stageRef.value.offsetWidth
  stageConfig.value.height = stageRef.value.offsetHeight
  cropConfig.value.width = stageConfig.value.width
  cropConfig.value.height = stageConfig.value.height
  videoSizeConfig.value.offsetWidth = mediaBox.value.offsetWidth
  videoSizeConfig.value.offsetHeight = mediaBox.value.offsetHeight
  videoSizeConfig.value.videoWidth = mediaBox.value.videoWidth
  videoSizeConfig.value.videoHeight = mediaBox.value.videoHeight
  videoSizeConfig.value.scaleX = videoSizeConfig.value.videoWidth / videoSizeConfig.value.offsetWidth
  videoSizeConfig.value.scaleY = videoSizeConfig.value.videoHeight / videoSizeConfig.value.offsetHeight
  //const parentElement = stageRef.value.parentElement
  var cropNode = cropRef.value.getNode();
  let node=cropRef.value.getStage()
  var trNode = trRef.value.getNode();
  trNode.nodes([cropNode]);
  trNode.getLayer().batchDraw();
  nextTick(()=> {
    node.fire('transformend', { customData: '這是模擬的資料' }, true);
  })
}

function timeRangeInput(e) {
  if(e[0]!=lastTimeRange.value[0]) mediaBox.value.currentTime=e[0]
  else if(e[1]!=lastTimeRange.value[1]) mediaBox.value.currentTime=e[1]
  lastTimeRange.value=e
}

function onLoaded(e) {
    console.log("影片載入成功",e)
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

function updateCropDataConfig(e) {
  cropDataConfig.value.x=e.target.attrs.x
  cropDataConfig.value.y=e.target.attrs.y
  cropDataConfig.value.width=e.target.attrs.width * (e.target.attrs?.scaleX || 1)
  cropDataConfig.value.height=e.target.attrs.height * (e.target.attrs?.scaleY || 1)
}

async function exportFunc() {
  let data= {
    info:props.info,
    controlConfig:controlConfig.value,
    cropDataConfig:cropDataConfig.value,
    timeRange:timeRange.value,
    videoSizeConfig:videoSizeConfig.value,
    exportConfig:exportConfig.value
  }

  let res=await ExportMedia(data)
}

async function testC() {
  let node=cropRef.value.getStage()
  node.fire('transformend', { customData: '這是模擬的資料' }, true);
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