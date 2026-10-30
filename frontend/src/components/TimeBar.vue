<template>
  <div class="timeBar">
    <div class="timeBarBox bg-light bg-gradient" ref="timeBarBox" >
        <div class="timeBarBoxMask" :style="{'left':timeBarBoxMaskData.left+'px','width':timeBarBoxMaskData.width+'px'}" v-show="controlPoint1isDown || controlPoint2isDown || controlPoint3isDown"></div>
        <div class="timeBarBoxCurrent"></div>
        <div class="timeBarBoxMark"></div>
        <div class="timeBarBoxControl">
            <div class="timeBarBoxControlRange bg-danger bg-gradient" :style="{'left':rangeLeft+'px','width':rangeWidth+'px'}"></div>
            <div class="timeBarBoxControlPoint point1" @mousemove.stop.prevent :style="{'left':controlPoint1+'px'}" @mousedown="controlPoint1isDown=true">
              <div class="timeBarBoxControlPointLabel">{{ xToTime(controlPoint1) }}</div>
            </div>
            <div class="timeBarBoxControlPoint point2" @mousemove.stop.prevent :style="{'left':controlPoint2+'px'}" @mousedown="controlPoint2isDown=true">
              <div class="timeBarBoxControlPointLabel">{{ xToTime(controlPoint2) }}</div>
            </div>
            <div class="timeBarBoxControlPoint point3" @mousemove.stop.prevent :style="{'left':controlPoint3+'px'}" @mousedown="controlPoint3isDown=true">
              <div class="timeBarBoxControlPointLabel">{{ xToTime(controlPoint3) }}</div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import {ref,computed,onMounted,nextTick} from "vue"

const props = defineProps(['min','max','step'])
const timeBarBox = ref(null)

const timeBarBoxMaskData = ref({
  left:0,
  width:0
})


const controlPoint1isDown = ref(false)
const controlPoint1 = ref(0)

const controlPoint2isDown = ref(false)
const controlPoint2 = ref(0)

const controlPoint3isDown = ref(false)
const controlPoint3 = ref(0)

const rangeLeft=computed(()=> {
  let res=controlPoint1.value
  if(controlPoint2.value < controlPoint1.value) res=controlPoint2.value
  return res
})

const rangeWidth=computed(()=> {
  let res=controlPoint2.value - controlPoint1.value
  if(controlPoint2.value < controlPoint1.value) res=controlPoint1.value - controlPoint2.value
  return res
})

onMounted(()=> {
  window.addEventListener("resize",()=> {
    nextTick(()=> {
      getNodeList()
      setTimeBarBoxMaskData()
    })
  })
  window.addEventListener("mousemove",(e)=> {
    move(e)
  })
  window.addEventListener("mouseup",(e)=> {
    controlPoint1isDown.value=false
    controlPoint2isDown.value=false
    controlPoint3isDown.value=false
  })
  //console.log(timeBarBox.value.offsetWidth)
  //console.log(nodeList.value)
  getNodeList()
  setTimeBarBoxMaskData()
})

function getDecimalPlaces(num) {
  let numStr = num.toString();
  if (!numStr.includes('.')) {
    return 0;
  }
  return numStr.split('.')[1].length;
}

function move(e) {
  if(e.target.className=="timeBarBoxMask") {
    if(controlPoint1isDown.value) {
      console.log(e.offsetX,e)
      controlPoint1.value = e.offsetX
      if(controlPoint1.value > timeBarBox.value.offsetWidth) controlPoint1.value=timeBarBox.value.offsetWidth
      if(controlPoint1.value < 0) controlPoint1.value=0
    }
    if(controlPoint2isDown.value) {
      console.log(e.offsetX,e)
      controlPoint2.value = e.offsetX
      if(controlPoint2.value > timeBarBox.value.offsetWidth) controlPoint2.value=timeBarBox.value.offsetWidth
      if(controlPoint2.value < 0) controlPoint2.value=0
    }
    if(controlPoint3isDown.value) {
      console.log(e.offsetX,e)
      controlPoint3.value = e.offsetX
      if(controlPoint3.value > timeBarBox.value.offsetWidth) controlPoint3.value=timeBarBox.value.offsetWidth
      if(controlPoint3.value < 0) controlPoint3.value=0
    }
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

function setTimeBarBoxMaskData() {
  timeBarBoxMaskData.value.left = timeBarBox.value.offsetLeft
  timeBarBoxMaskData.value.width = timeBarBox.value.offsetWidth
}
</script>