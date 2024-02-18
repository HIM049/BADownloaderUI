<template>
    <el-header style="--wails-draggable:drag; margin-bottom:70px">
        <nav style=" display:flex; justify-content:space-between; margin: 20px 0px 20px 15px;">
            <div id="logo" style=" display:flex; align-items:center; ">
                <svg style="width: 40px;" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 396.65 396"><defs><clipPath id="a" transform="translate(-747 -329)"><rect x="850.03" y="452.03" width="174.85" height="176.95" style="fill:none"/></clipPath></defs><path d="M1089.72,430a8.3,8.3,0,0,1,8.28,8.28V642.72a8.3,8.3,0,0,1-8.28,8.28H785.28a8.3,8.3,0,0,1-8.28-8.28V438.28a8.3,8.3,0,0,1,8.28-8.28h304.44m0-30H785.28A38.28,38.28,0,0,0,747,438.28V642.72A38.28,38.28,0,0,0,785.28,681h304.44A38.28,38.28,0,0,0,1128,642.72V438.28A38.28,38.28,0,0,0,1089.72,400Z" transform="translate(-747 -329)" style="fill:#ffc7c7"/><line x1="128.5" y1="78.5" x2="114" y2="15" style="fill:none;stroke:#ffc7c7;stroke-linecap:round;stroke-miterlimit:10;stroke-width:30px"/><line x1="249.5" y1="78.5" x2="267" y2="15" style="fill:none;stroke:#ffc7c7;stroke-linecap:round;stroke-miterlimit:10;stroke-width:30px"/><polyline points="331 213 330.37 361.58 331 386" style="fill:none;stroke:#333;stroke-linecap:round;stroke-linejoin:round;stroke-width:20px"/><polyline points="386.65 319.24 331 386 270.29 319.24" style="fill:none;stroke:#333;stroke-linecap:round;stroke-linejoin:round;stroke-width:20px"/><g style="clip-path:url(#a)"><path d="M937.39,452.17c3.71,0,6.72,3.3,6.72,7.36V621.36c0,4.06-3,7.35-6.72,7.35s-6.72-3.29-6.72-7.35V459.53c0-4.06,3-7.36,6.72-7.36ZM897.05,481.6c3.72,0,6.73,3.29,6.73,7.35v103c0,4.07-3,7.36-6.73,7.36s-6.72-3.29-6.72-7.36V489c0-4.06,3-7.35,6.72-7.35Zm80.67,0c3.71,0,6.72,3.29,6.72,7.35v103c0,4.07-3,7.36-6.72,7.36S971,596,971,591.93V489c0-4.06,3-7.35,6.72-7.35ZM1018.05,511c3.72,0,6.73,3.29,6.73,7.35v44.14c0,4.06-3,7.36-6.73,7.36s-6.72-3.3-6.72-7.36V518.37c0-4.06,3-7.35,6.72-7.35Zm-161.33,0c3.71,0,6.72,3.29,6.72,7.35v44.14c0,4.06-3,7.36-6.72,7.36s-6.72-3.3-6.72-7.36V518.37c0-4.06,3-7.35,6.72-7.35Zm0,0" transform="translate(-747 -329)" style="fill:#8a8a8a"/></g></svg>
                <div style="display:flex; font-size:20px; font-weight:800; margin-left:10px;">
                    <p style="color: #ffabab;">B</p>
                    <p>ili</p>
                    <p style="color: #ffabab;">A</p>
                    <p>udio</p>
                    <p style="color: #ffabab;">D</p>
                    <p>ownloader</p>
                </div>
            </div>
            <div id="ctl-buttons">
                <button class="ctl-button" id="minimize" @click="WindowMinimise()"><svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 -960 960 960"><path d="M200-440v-80h560v80H200Z"/></svg></button>
                <button class="ctl-button" id="close" @click="Quit()"><svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 -960 960 960"><path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z"/></svg></button>
            </div>
        </nav>
        <div id="navbar" style=" display: flex; align-items: center; justify-content: center;">
            <div class="nav-button" :data-myattribute="model == -12 ? 'current' : ''">
                <button class="button" id="space" @click="model = -12">
                    <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" fill="currentColor" class="bi bi-file-person" viewBox="0 0 16 16"><path d="M12 1a1 1 0 0 1 1 1v10.755S12 11 8 11s-5 1.755-5 1.755V2a1 1 0 0 1 1-1h8zM4 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H4z"/><path d="M8 10a3 3 0 1 0 0-6 3 3 0 0 0 0 6z"/></svg>
                </button>
                <label for="space">个人空间</label>
            </div>

            <div class="nav-button" :data-myattribute="model >= 0 ? 'current' : ''">
                <button class="button" id="download" @click="model = 0">
                    <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" fill="currentColor" class="bi bi-file-earmark-arrow-down" viewBox="0 0 16 16"><path d="M8.5 6.5a.5.5 0 0 0-1 0v3.793L6.354 9.146a.5.5 0 1 0-.708.708l2 2a.5.5 0 0 0 .708 0l2-2a.5.5 0 0 0-.708-.708L8.5 10.293V6.5z"/><path d="M14 14V4.5L9.5 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2zM9.5 3A1.5 1.5 0 0 0 11 4.5h2V14a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1h5.5v2z"/></svg>
                </button>
                <label for="download">批量下载</label>
            </div>

            <div class="nav-button" :data-myattribute="model == -11 ? 'current' : ''">
                <button class="button" id="single" @click="model = -11">
                    <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" fill="currentColor" class="bi bi-file-earmark-music" viewBox="0 0 16 16"><path d="M11 6.64a1 1 0 0 0-1.243-.97l-1 .25A1 1 0 0 0 8 6.89v4.306A2.572 2.572 0 0 0 7 11c-.5 0-.974.134-1.338.377-.36.24-.662.628-.662 1.123s.301.883.662 1.123c.364.243.839.377 1.338.377.5 0 .974-.134 1.338-.377.36-.24.662-.628.662-1.123V8.89l2-.5V6.64z"/><path d="M14 14V4.5L9.5 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2zM9.5 3A1.5 1.5 0 0 0 11 4.5h2V14a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1h5.5v2z"/></svg>
                </button>
                <label for="single">单曲下载</label>
            </div>

            <div class="nav-button" :data-myattribute="model == -10 ? 'current' : ''">
                <button class="button" id="settings" @click="model = -10">
                    <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" fill="currentColor" class="bi bi-gear" viewBox="0 0 16 16"><path d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492zM5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0z"/><path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52l-.094-.319zm-2.633.283c.246-.835 1.428-.835 1.674 0l.094.319a1.873 1.873 0 0 0 2.693 1.115l.291-.16c.764-.415 1.6.42 1.184 1.185l-.159.292a1.873 1.873 0 0 0 1.116 2.692l.318.094c.835.246.835 1.428 0 1.674l-.319.094a1.873 1.873 0 0 0-1.115 2.693l.16.291c.415.764-.42 1.6-1.185 1.184l-.291-.159a1.873 1.873 0 0 0-2.693 1.116l-.094.318c-.246.835-1.428.835-1.674 0l-.094-.319a1.873 1.873 0 0 0-2.692-1.115l-.292.16c-.764.415-1.6-.42-1.184-1.185l.159-.291A1.873 1.873 0 0 0 1.945 8.93l-.319-.094c-.835-.246-.835-1.428 0-1.674l.319-.094A1.873 1.873 0 0 0 3.06 4.377l-.16-.292c-.415-.764.42-1.6 1.185-1.184l.292.159a1.873 1.873 0 0 0 2.692-1.115l.094-.319z"/></svg>
                </button>
                <label for="settings">设置</label>
            </div>

        </div>
    </el-header>
</template>

<script setup>
import { WindowMinimise, Quit } from '../../../wailsjs/runtime'
import { ref, reactive, computed, defineProps, defineModel } from 'vue'

const model = defineModel()

</script>

<style>

/* 顶栏窗口操作按钮 */
nav #ctl-buttons button.ctl-button svg {
  height: 25px;
  width: 25px;
  padding: 5px;
  margin-left: 5px;
  border-radius: 4px;
}

/* 鼠标选中时将图标改为白色 */
nav #ctl-buttons button.ctl-button:hover svg path {
  fill: #fff;
}
/* 退出按钮选中背景色 */
nav div #close:hover svg {
  background-color: #e94c4c;
}
/* 最小化按钮选中背景色 */
nav div #minimize:hover svg {
  background-color: #acacac;
}

/* 导航栏按钮 */
header #navbar div.nav-button {
  margin-left: 12px;
  display: flex;
  align-items: center;
  border-bottom: 2px rgba(0, 0, 0, 0) solid;
  transition: 0.2s;
}

/* 图标按钮去除默认样式 */
button {
  background: none;
  border: none;
}

/* 按钮底部指示条 */
header #navbar div.nav-button[data-myattribute="current"] {
  border-bottom: 2px #cf8888 solid;
}

/* 按钮标题与图标颜色 */
/* 鼠标选中时 */
header #navbar div.nav-button:hover button.button svg path,
header #navbar div.nav-button:hover label,
/* 当前页面 */
header #navbar div.nav-button[data-myattribute="current"] button.button svg path, 
header #navbar div.nav-button[data-myattribute="current"] label {
  color: #cf8888;
  fill: #cf8888;
}
</style>