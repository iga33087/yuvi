<template>
  <div class="home">
    <div class="fs-2 mb-4">請選擇本地媒體或輸入影片URL</div>
    <CardBox class="rounded-5">
        <template v-slot:content>
            <div class="d-flex align-items-center justify-content-between">
              <div class="fs-4">選擇本地媒體</div>
              <div class="d-flex align-items-center">
                <button type="button" class="btn btn-primary" @click="mediaChooseChange">選擇檔案</button>
                <button type="button" class="btn btn-primary" @click="mediaChooseChange2">選擇檔案2</button>
              </div>
            </div>
        </template>
    </CardBox>
    <CardBox class="rounded-5">
        <template v-slot:content>
            <div class="d-flex align-items-center justify-content-between">
              <div class="fs-4">輸入影片URL</div>
              <div class="d-flex align-items-center">
                <input type="text" class="form-control me-2" style="width:400px;" placeholder="影片URL">
                <button type="button" class="btn btn-primary">提交</button>
              </div>
            </div>
        </template>
    </CardBox>

  </div>
</template>

<script setup>
import { Test2,SelectAndProcessFile,GetFileInfo } from "../../wailsjs/go/main/App";
import CardBox from '../components/CardBox.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

async function mediaChooseChange(e) {
    let path=await SelectAndProcessFile()
    router.push({name:"Result",query:{path}})
}
async function mediaChooseChange2(e) {
    let path=await SelectAndProcessFile()
    router.push({name:"Edit",query:{path}})
}
</script>