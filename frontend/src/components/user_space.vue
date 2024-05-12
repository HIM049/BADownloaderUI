<template>
    <FramePage title="个人空间">
        <var-space style="align-items: center;" justify="center" v-if="!showLoginWindow && !is_login">
            你还没有登录
            <var-button type="primary" @click="login">登录</var-button>
        </var-space>
        <var-collapse-transition :expand="showLoginWindow">
            <var-paper style="background-color: var(--color-primary-container);">
                <h3>扫码登录</h3>
                <img style="display: flex; margin: 0 auto; border-radius: 15px;" :src="qrcodeStr" alt="Image">
                <p style="display: flex; justify-content: center;">{{ loginText }}</p>
            </var-paper>
        </var-collapse-transition>

        <var-paper style="background-color: var(--color-primary-container); margin-bottom: 20px;" v-if="is_login">
            <h3>创建的收藏夹</h3>
            <var-cell v-for="(collect, index) in user_collect" style="margin: 0 10px;">
                {{ collect.title }}
                <template #extra>
                    <var-button type="primary" @click="">下载</var-button>
                </template>
            </var-cell>
        </var-paper>

        <var-paper style="background-color: var(--color-primary-container);" v-if="is_login">
            <h3>收藏和订阅</h3>
            <var-cell v-for="(collect, index) in user_Favourite" style="margin: 0 10px;">
                {{ collect.title }}
                <template #extra>
                    <var-button type="primary" @click="">下载</var-button>
                </template>
            </var-cell>
            <var-space style="display: flex; align-items: center;">
                <var-button-group type="primary" size="normal" outline >
                    <var-button @click="page_index--">上一页</var-button>
                    <var-button @click="page_index++">下一页</var-button>
                </var-button-group>
            </var-space>
        </var-paper>
        
    </FramePage>
</template>

<script setup>
import FramePage from '../components/modules/frame_page.vue'
import { ref, reactive, onMounted, watch } from 'vue'
import { LoginBilibili, LoadConfig, GetUsersCollect, GetFavCollect } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOnce } from '../../wailsjs/runtime'
import { Snackbar } from '@varlet/ui'

const loginText = ref("请扫描二维码登录")
const qrcodeStr = ref(null)
const showLoginWindow = ref(false)
const is_login = ref(false)
const page_index = ref(1)

const user_collect = ref([]);
const user_Favourite = ref([]);
const fav_count = ref(0);

onMounted(() => {
    checkLogin();
    if (is_login) {
        GetUsersCollect().then(result => {
            user_collect.value = result.List;
        })
        getFavCollect();
    }
}),

watch(page_index, (newValue) => {
    if (newValue <= 0) {
        page_index.value = 1;
    }else if (newValue > fav_count.value / 20 + 1) {
        page_index.value = parseInt(fav_count.value / 20 + 1);
    }
    getFavCollect()

})

EventsOn("qrcodeStr", (qr)=>{
    qrcodeStr.value = "data:image/png;base64," + qr;
})

EventsOn("loginStatus", (status) => {
    if (status == "登录成功") {
        Snackbar.success("登录成功");
        checkLogin();
        setTimeout(() => {
            window.location.reload();
        }, 1000);
    }
    if (status == "二维码已失效") {
        Snackbar.warning("二维码已失效");
        showLoginWindow.value = false;
    }
    loginText.value = status;
})

function getFavCollect() {
    GetFavCollect(page_index.value).then(result => {
        user_Favourite.value = result.List
        fav_count.value = result.count
    })
}

function checkLogin() {
    LoadConfig().then(result => {
        is_login.value = result.Account.is_login
    })
}

function login() {
    LoginBilibili().then(result => {
        showLoginWindow.value = false;
    })
    showLoginWindow.value = true;
}
</script>