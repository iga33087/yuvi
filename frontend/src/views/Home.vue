<template>
  <div class="home">
    <div class="fs-2 mb-4">請選擇本地媒體或輸入影片URL</div>
    <CardBox class="rounded-5">
        <template v-slot:content>
            <div class="d-flex align-items-center justify-content-between">
              <div class="d-flex align-items-center fs-4"><HostIcon class="me-4" />選擇本地媒體</div>
              <div class="d-flex align-items-center">
                <button type="button" class="btn d-flex align-items-center btn-primary bg-gradient" @click="mediaChooseChange"><FileIcon class="me-2" />選擇檔案</button>
              </div>
            </div>
        </template>
    </CardBox>
    <CardBox class="rounded-5">
        <template v-slot:content>
            <div class="d-flex align-items-center justify-content-between">
              <div class="d-flex align-items-center fs-4"><GlobeIcon class="me-4" />輸入影片URL</div>
              <div class="d-flex align-items-center">
                <input type="text" class="form-control me-2" style="width:400px;" placeholder="影片URL">
                <button type="button" class="btn btn-primary bg-gradient">提交</button>
              </div>
            </div>
        </template>
    </CardBox>

    <TimeBar :min="0" :max="10" :step="0.1" />

  </div>
</template>

<script setup>
import FileIcon from "../assets/icons/file-earmark-play.svg"
import HostIcon from "../assets/icons/pc-display-horizontal.svg"
import GlobeIcon from "../assets/icons/globe2.svg"
import { Test2,SelectAndProcessFile,GetFileInfo } from "../../wailsjs/go/main/App";
import CardBox from '../components/CardBox.vue'
import TimeBar from '../components/TimeBar.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

async function mediaChooseChange(e) {
    let path=await SelectAndProcessFile()
    router.push({name:"Edit",query:{path}})
}
</script>