<script setup>
import { reactive, computed, ref } from 'vue'
import { SearchFavListInformation, MakeAndSaveList, StartDownload , MakeUpEditor } from '../../wailsjs/go/main/App'
import { DocumentAdd, Edit, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus';

const parms = reactive({
  favListID: "",
})

// 查询函数返回值
const resp = reactive({
  title: "",
  cover: "",
  count: 0,
  up_name: "",
  up_avatar: "",
})

// 下载设置
const options = reactive({
  downCount: 0,
  downPart: true,
  songName: true,
  songCover: true,
  songAuthor: true,
})

// 组件显示状态
const status = reactive({
  favListInf: false,
  makeingList: false,
  downloading: false,
  makeListButton: false,
  editlistButton: false,
  downloadButton: false,
})

// const downFileType = ref("mp3")

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
    console.log(result);
    if (result.message == "0") {
      resp.title = result.Data.Info.title
      resp.cover = result.Data.Info.cover
      resp.count = result.Data.Info.media_count
      resp.up_name = result.Data.Info.Upper.name
      resp.up_avatar = result.Data.Info.Upper.face
      
      // 开放创建列表按钮
      status.makeListButton = true
    } else {
      // 无效的收藏夹
      resp.title = "无效的收藏夹"
      ElMessage.warning("无效的收藏夹")
      // 关闭创建列表按钮
      status.makeListButton = false
    }
    status.favListInf = true
  })
  // 展示收藏夹信息
}

// 生成视频列表
function creatVideoList() {
  ElMessage.success("正在创建视频列表");
  // 关闭创建列表按钮
  status.makeListButton = false

  status.makeingList = true;
  MakeAndSaveList(parms.favListID, Number(options.downCount), options.downPart).then(result => {
    if (result != null) {
      // 创建失败
      ElMessage.error(result);
    } else {
      // 创建成功
      ElMessage.success("创建完成");
    }
    
    // 修改组件状态
    status.makeingList = false;
    status.makeListButton = true
    status.editlistButton = true;
    status.downloadButton = true;
  })
}

function downloadVideoList() {
  ElMessage.success("开始下载");
  // 关闭下载按钮
  status.downloadButton = false

  status.downloading = true;

  var opt = {
    song_name: options.songName,
    song_cover: options.songCover,
    song_author: options.songAuthor,
  }
  StartDownload(opt).then(result => {
    status.downloading = false;
    
    status.downloadButton = true
    ElMessage.success("下载完成");
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
      <input type="text" class="input" id="favIdInput" v-model.trim="parms.favListID" @input.lazy="queryFavListInformation"
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
  <el-form class="card" label-position="right" label-width="100px">
    <h3>下载选项</h3>
    
    <el-form-item label="下载数量">
      <el-input-number v-model="options.downCount" :min="0" :max="resp.count" />
    </el-form-item>

    <el-form-item label="下载分集">
      <el-switch v-model="options.downPart" style="--el-switch-on-color: #13ce66;" />
    </el-form-item>

  </el-form>

  <el-form class="card" label-position="right" label-width="100px">
    <h3>元数据选项</h3>

    <el-form-item label="歌曲名称">
      <el-switch v-model="options.songName" style="--el-switch-on-color: #13ce66;" />
    </el-form-item>
    <el-form-item label="歌曲封面">
      <el-switch v-model="options.songCover" style="--el-switch-on-color: #13ce66;" />
    </el-form-item>
    <el-form-item label="歌曲作者">
      <el-switch v-model="options.songAuthor" style="--el-switch-on-color: #13ce66;" />
    </el-form-item>

  </el-form>

  <!-- <form class="card">
    <h3>文件选项</h3>
    
    <label for="file-2">MP3</label>
    <input type="radio" id="file-2" v-model="downFileType" value="mp3" />
    <label for="file-1">M4A</label>
    <input type="radio" id="file-1" v-model="downFileType" value="m4a" />
  </form> -->

  <el-form>
    <el-button class="button2" type="info" size="large" :icon="DocumentAdd" plain @click="creatVideoList" :loading="status.makeingList" :disabled="!status.makeListButton">生成视频列表</el-button><br/>
    <el-button class="button2" type="info" size="large" :icon="Edit" plain @click="makeUpEditor" :disabled="!status.editlistButton">编辑列表信息</el-button><br/>
    <el-button class="button2" type="success" size="large" :icon="Download" plain @click="downloadVideoList" :disabled="!status.downloadButton" :loading="status.downloading">下载列表内容</el-button><br/>
  </el-form>
</template>

<style>

.input {
    /* 收藏夹信息表单输入框 */
    width: 100%;
    height: 30px;
    text-align: center;
    border: 2px solid var(--el-border-color-darker);
    border-radius: var(--el-border-radius-base);
    padding: 2px 0;
}

.button2 {
    margin-bottom: 10px;
    width: 100%;
    height: 45px;
    font-size: 20px;
    color: #000;
}

/* 样式修改 */
.el-form-item__content {
    justify-content: center;
}
.el-form-item {
    margin-bottom: 10px;
}
</style>