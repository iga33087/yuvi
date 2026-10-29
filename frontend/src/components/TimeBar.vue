<template>
  <div class="timeBar">
    <div class="timeBarBox" ref="timeBarBox" >
        <div class="timeBarBoxCurrent"></div>
        <div class="timeBarBoxMark"></div>
        <div class="timeBarBoxControl">
            <div class="timeBarBoxControlRange"></div>
            <div class="timeBarBoxControlPoint" @mousemove.stop.prevent :style="{'left':controlPoint1+'px'}" @mousedown="controlPoint1isDown=true">
              <div class="timeBarBoxControlPointLabel">{{ xToTime(controlPoint1) }}</div>
            </div>
            <div class=""></div>
        </div>
    </div>
  </div>
</template>

<script setup>
import {ref,computed,onMounted,nextTick} from "vue"

const props = defineProps(['min','max','step'])
const timeBarBox = ref(null)

const controlPoint1isDown = ref(false)
const controlPoint1 = ref(0)

onMounted(()=> {
  window.addEventListener("resize",()=> {
    nextTick(()=> {
      getNodeList()
    })
  })
  window.addEventListener("mousemove",(e)=> {
    move(e)
  })
  window.addEventListener("mouseup",(e)=> {
    controlPoint1isDown.value=false
  })
  //console.log(timeBarBox.value.offsetWidth)
  //console.log(nodeList.value)
  getNodeList()
})

function getDecimalPlaces(num) {
  let numStr = num.toString();
  if (!numStr.includes('.')) {
    return 0;
  }
  return numStr.split('.')[1].length;
}

function move(e) {
  if(controlPoint1isDown.value) {
    console.log(e.offsetX,e)
    controlPoint1.value += e.movementX
    if(controlPoint1.value > timeBarBox.value.offsetWidth) controlPoint1.value=timeBarBox.value.offsetWidth
    if(controlPoint1.value < 0) controlPoint1.value=0
  }
}

function getNodeList() {
  let res = {}
  let fixed = getDecimalPlaces(props.step)
  let range = (props.max - props.min) / props.step
  let each = timeBarBox.value.offsetWidth / range
  let index = 0
  let value = 0
  while(value <= timeBarBox.value.offsetWidth) {
    res[index] = value
    index = Number((index+=props.step).toFixed(fixed))
    value = Number((value+=each).toFixed(0))
  }
  console.log("getNodeList",res)
  return res
}

function xToTime(x) {
  if(timeBarBox.value) {
    let fixed = getDecimalPlaces(props.step)
    let range = (props.max - props.min) / props.step
    let each = range / timeBarBox.value.offsetWidth
    return Number(((x*each)*props.step).toFixed(fixed))
  }
}
</script>