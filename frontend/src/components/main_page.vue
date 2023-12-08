<script setup>
import { reactive, computed, ref } from 'vue'
import { SearchFavListInformation, MakeAndSaveList, StartDownload , MakeUpEditor } from '../../wailsjs/go/main/App'
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

const options = reactive({
  downCount: "0",
  downPart: true,
  songName: true,
  songCover: true,
  songAuthor: true,
})
const status = reactive({
  favListInf: false,
  makeingList: false,
  downloading: false,
  editlistButton: false,
  downloadButton: false,
})
const downFileType = ref("mp3")


// 742380048

// 查询收藏夹信息
function queryFavListInformation() {
  // TODO: 输入校验 / 直接输入 http 链接
  parms.favListID = parms.favListID.replace(/\D/g, '');
  if (parms.favListID.length > 18){
    parms.favListID = parms.favListID.slice(0, 18);
  }
  SearchFavListInformation(parms.favListID).then(result => {
    // 判断信息有效性
    if (result.message == "0") {
      resp.title = result.Data.Info.title
      resp.cover = result.Data.Info.cover
      resp.count = result.Data.Info.media_count
      resp.up_name = result.Data.Info.Upper.name
      resp.up_avatar = result.Data.Info.Upper.face
    } else {
      // 无效的收藏夹
      console.error("无效的收藏夹");
      resp.title = "无效的收藏夹"
      // TODO：制作收藏夹无效警告框
    }
    status.favListInf = true
  })
  // 展示收藏夹信息
}

// 生成视频列表
function creatVideoList() {
  console.log("正在创建视频列表");
  status.makeingList = true;
  MakeAndSaveList(parms.favListID, Number(downCount.value), options.downPart).then(result => {
    if (result != null) {
      console.error(result);
    } else {
      console.log("OK!");
    }
    status.makeingList = false;
    status.editlistButton = true;
    status.downloadButton = true;
  })
}

function downloadVideoList() {
  console.log("正在下载");
  status.downloading = true;

  var opt = {
    song_name: options.songName,
    song_cover: options.songCover,
    song_author: options.songAuthor,
  }
  StartDownload(opt).then(result => {
    status.downloading = false;
  })
}

// 启动文本编辑器（临时解决方案）
function makeUpEditor() {
  MakeUpEditor()
}

</script>

<template>
  <!-- 收藏夹信息输入及展示 -->
  <div id="fav-input" class="fav-input">
    <form>
      <input type="text" id="favIdInput" v-model.trim="parms.favListID" @input.lazy="queryFavListInformation"
        placeholder="请输入要下载的收藏夹 ID">
    </form>

    <div class="information-box card" v-show="status.favListInf">
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
    <li class="option-list">
        <label for="downCount">下载数量</label>
        <input type="text" id="downCount" v-model="options.downCount" style="width: 40px;">
    </li>
      <li class="option-list">
        <label for="dPart">下载视频分 P</label>
        <input type="checkbox" class="switch" id="dPart" v-model="options.downPart">
        <label for="dPart" class="switch-label green"></label>
      </li>
    </ol>
  </form>

  <form id="input-meta" class="card">
    <h3>元数据选项</h3>
    <ol>
      <li class="option-list">
        <label for="songName">歌曲名称</label>
        <input type="checkbox" class="switch" id="songName" value="songName" v-model="options.songName">
        <label for="songName" class="switch-label green"></label>
      </li>

      <li class="option-list">
        <label for="songCover">歌曲封面</label>
        <input type="checkbox" class="switch" id="songCover" value="songCover" v-model="options.songCover">
        <label for="songCover" class="switch-label green"></label>
      </li>

      <li class="option-list">
        <label for="songAuthor">歌曲作者</label>
        <input type="checkbox" class="switch" id="songAuthor" value="songAuthor" v-model="options.songAuthor">
        <label for="songAuthor" class="switch-label green"></label>
      </li>
    </ol>
  </form>

  <!-- <form class="card">
    <h3>文件选项</h3>
    
    <label for="file-2">MP3</label>
    <input type="radio" id="file-2" v-model="downFileType" value="mp3" />
    <label for="file-1">M4A</label>
    <input type="radio" id="file-1" v-model="downFileType" value="m4a" />
  </form> -->

  <form>
    <button class="button" @click.prevent="creatVideoList">
      <p>生成视频列表</p>
      <loadingAnimation v-if="status.makeingList"/>
      <svg v-else t="1700970936750" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="15541" width="25" height="25"><path d="M358.165868 554.624796c-0.533143 0.680499-1.066285 1.391696-1.303692 2.251274l-41.104163 150.700257c-2.400676 8.772804 0.059352 18.226107 6.550183 24.892947 4.860704 4.742001 11.261485 7.350408 18.077727 7.350408 2.252297 0 4.504594-0.267083 6.727215-0.860601l149.630902-40.808428c0.23843 0 0.357134 0.207731 0.534166 0.207731 1.718131 0 3.408633-0.62217 4.683672-1.927909l400.113747-400.054395c11.883655-11.897981 18.404162-28.109198 18.404162-45.74281 0-19.989263-8.476045-39.963177-23.324218-54.767348l-37.786605-37.844933c-14.81645-14.848173-34.822087-23.338544-54.797024-23.338544-17.631566 0-33.842783 6.520507-45.75816 18.388812L358.758362 553.232077C358.344946 553.615816 358.462626 554.179658 358.165868 554.624796M862.924953 257.19778l-39.742143 39.71349-64.428382-65.451688 39.180348-39.179324c6.193049-6.222725 18.194384-5.318122 25.308409 1.822508l37.813211 37.845956c3.943822 3.941775 6.195096 9.18622 6.195096 14.372336C867.223862 250.574942 865.710392 254.42769 862.924953 257.19778M429.322487 560.907896l288.712541-288.728914 64.459081 65.49569L494.314711 625.838721 429.322487 560.907896zM376.718409 677.970032l20.863167-76.580143 55.656601 55.657624L376.718409 677.970032z" fill="#231F20" p-id="15542"></path><path d="M888.265084 415.735539c-15.144932 0-27.562752 12.313443-27.620058 27.665083l0 372.98283c0 19.559475-15.885805 35.444257-35.475979 35.444257L194.220958 851.827709c-19.559475 0-35.504632-15.883759-35.504632-35.444257L158.716326 207.602222c0-19.575848 15.945157-35.474956 35.504632-35.474956l406.367171 0c15.231913 0 27.592428-12.371772 27.592428-27.606755 0-15.202237-12.360516-27.590382-27.592428-27.590382L190.013123 116.930129c-47.684022 0-86.49291 38.779212-86.49291 86.49291L103.520213 820.59231c0 47.713698 38.808888 86.47756 86.49291 86.47756l639.334083 0c47.715745 0 86.509283-38.763862 86.509283-86.47756L915.856489 443.222567C915.794068 428.048983 903.408993 415.735539 888.265084 415.735539" fill="#231F20" p-id="15543"></path></svg>
    </button>
    
    <button class="button" @click.prevent="makeUpEditor" :disabled="!status.editlistButton">
      <p>编辑列表信息</p>
      <svg t="1700971089790" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="15696" width="25" height="25"><path d="M387.657143 906.971429L146.285714 738.742857 636.342857 36.571429 877.714286 204.8l-490.057143 702.171429zM248.685714 716.8L373.028571 804.571429l402.285715-577.828572-124.342857-87.771428-402.285715 577.828571z" fill="" p-id="15697"></path><path d="M160.914286 987.428571l65.828571-43.885714-7.314286-204.8H146.285714z" fill="" p-id="15698"></path><path d="M160.914286 987.428571l7.314285-80.457142 197.485715-65.828572 21.942857 65.828572zM753.371429 380.342857L533.942857 219.428571l43.885714-58.514285 219.428572 153.6z" fill="" p-id="15699"></path></svg>
    </button>

    <button class="button" @click.prevent="downloadVideoList" :disabled="!status.downloadButton">
      <p>开始下载</p>
      <loadingAnimation v-if="status.downloading"/>
      <svg v-else t="1700936201572" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="12672" width="25" height="25"><path d="M828.975746 894.125047 190.189132 894.125047c-70.550823 0-127.753639-57.18542-127.753639-127.752616L62.435493 606.674243c0-17.634636 14.308891-31.933293 31.93227-31.933293l63.889099 0c17.634636 0 31.93227 14.298658 31.93227 31.933293l0 95.821369c0 35.282574 28.596292 63.877843 63.87682 63.877843L765.098927 766.373455c35.281551 0 63.87682-28.595268 63.87682-63.877843l0-95.821369c0-17.634636 14.298658-31.933293 31.943526-31.933293l63.877843 0c17.634636 0 31.933293 14.298658 31.933293 31.933293l0 159.699212C956.729385 836.939627 899.538849 894.125047 828.975746 894.125047L828.975746 894.125047zM249.938957 267.509636c12.921287-12.919241 33.884738-12.919241 46.807049 0l148.97087 148.971893L445.716876 94.89323c0-17.634636 14.300704-31.94762 31.933293-31.94762l63.875796 0c17.637706 0 31.945573 14.312984 31.945573 31.94762l0 321.588299 148.97087-148.971893c12.921287-12.919241 33.875528-12.919241 46.796816 0l46.814212 46.818305c12.921287 12.922311 12.921287 33.874505 0 46.807049L552.261471 624.930025c-1.140986 1.137916-21.664416 13.68365-42.315758 13.69286-20.87647 0.010233-41.878806-12.541641-43.020816-13.69286L203.121676 361.13499c-12.922311-12.933567-12.922311-33.884738 0-46.807049L249.938957 267.509636 249.938957 267.509636z" fill="#272636" p-id="12673"></path></svg>
    </button>
  </form>
</template>
