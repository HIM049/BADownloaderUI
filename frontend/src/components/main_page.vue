<script setup>
import { reactive, computed, ref } from 'vue'
import { SearchFavListInformation, MakeAndSaveList, DownloadList } from '../../wailsjs/go/main/App'

import loadingAnimation from './loading.vue'

const parms = reactive({
  favListID: "",
})
const resp = reactive({
  title: "",
  cover: "",
  count: 0,
  up_name: "",
  up_avatar: "",
})
const DownPart = ref(true)
const showFavInf = ref(false)
const metaSettings = ref(["songName", "songCover", "songAuthor"])
const tips = ref("")
const showTips = ref(false)


// 742380048

// 查询收藏夹信息
function searchFavListInformation() {
  if (parms.favListID.length > 2) {
    SearchFavListInformation(parms.favListID).then(result => {
      resp.title = result.Data.Info.title
      resp.cover = result.Data.Info.cover
      resp.count = result.Data.Info.media_count
      resp.up_name = result.Data.Info.Upper.name
      resp.up_avatar = result.Data.Info.Upper.face
    })
    // 展示收藏夹信息
    showFavInf.value = true
  }
}

// 生成视频列表
function creatVideoList() {
  showTips.value = true
  tips.value = "正在生成列表"
  MakeAndSaveList("C:/Users/HIM/Desktop/Download/video_list.json", parms.favListID, 0, false).then(result => {
    if (result != null) {
      LogError(result)
    } else {
      LogInfo("OK!")
    }
    showTips.value = false
  })
}

function downloadVideoList() {
  DownloadList("C:/Users/HIM/Desktop/Download/video_list.json", 1)
}

</script>

<template>
  <!-- 收藏夹信息输入及展示 -->
  <div id="fav-input" class="fav-input">
    <form>
      <input type="text" id="favIdInput" v-model.trim="parms.favListID" @input.lazy="searchFavListInformation"
        placeholder="请输入要下载的收藏夹 ID">
      <button @click.prevent="creatVideoList">
        <svg t="1700907443377" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4908" width="30" height="30"><path d="M385 840.5c-20.8 0-41.7-7.9-57.6-23.8L87.6 576.9c-31.8-31.8-31.8-83.3 0-115.1s83.3-31.8 115.1 0l239.8 239.8c31.8 31.8 31.8 83.3 0 115.1-15.9 15.9-36.7 23.8-57.5 23.8z" fill="#2c2c2c" p-id="4909"></path><path d="M384.6 840.5c-20.8 0-41.7-7.9-57.6-23.8-31.8-31.8-31.8-83.3 0-115.1l494.2-494.2c31.8-31.8 83.3-31.8 115.1 0s31.8 83.3 0 115.1L442.2 816.7c-15.9 15.9-36.8 23.8-57.6 23.8z" fill="#2c2c2c" p-id="4910"></path></svg>
      </button>
    </form>

    <div class="loading-tips" v-show="showTips">
      <p>{{ tips }}</p>
      <loadingAnimation />
    </div>

    <div class="information-box card" v-show="showFavInf">
      <h3>收藏夹信息</h3>
      <div class="favlist-box">
        <img :src="resp.cover">
        <div>
          <h3>{{ resp.title }}</h3>
          <h4>{{ resp.count }} 个内容</h4>
          <div class="up-box">
            <img :src="resp.up_avatar">
            <h4>{{ resp.up_name }}</h4>
          </div>
        </div>
      </div>
    </div>

  </div>

  <!-- 下载设定 -->
  <form class="card">
    <h3>下载选项</h3>
    <ol>
      <li class="checkbox">
        <label for="dPart">下载视频分 P</label>
        <input type="checkbox" class="switch" id="dPart" v-model="DownPart">
        <label for="dPart" class="switch-label green"></label>
      </li>
    </ol>
  </form>
  <form id="input-meta" class="card">
    <h3>元数据选项</h3>
    <ol>
      <li class="checkbox">
        <label for="songName">歌曲名称</label>
        <input type="checkbox" class="switch" id="songName" value="songName" v-model="metaSettings">
        <label for="songName" class="switch-label green"></label>
      </li>

      <li class="checkbox">
        <label for="songCover">歌曲封面</label>
        <input type="checkbox" class="switch" id="songCover" value="songCover" v-model="metaSettings">
        <label for="songCover" class="switch-label green"></label>
      </li>

      <li class="checkbox">
        <label for="songAuthor">歌曲作者</label>
        <input type="checkbox" class="switch" id="songAuthor" value="songAuthor" v-model="metaSettings">
        <label for="songAuthor" class="switch-label green"></label>
      </li>
    </ol>
  </form>
  <button @click="downloadVideoList">开始下载</button>
</template>
