<template>
    <FramePage title="个人空间">
        <var-space style="align-items: center;" justify="center" v-show="!showLoginWindow && !is_login">
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
        
    </FramePage>
</template>

<script setup>
import FramePage from '../components/modules/frame_page.vue'
import { ref, reactive, onMounted } from 'vue'
import { LoginBilibili, LoadConfig } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOnce } from '../../wailsjs/runtime'
import { Snackbar } from '@varlet/ui'

const loginText = ref("请扫描二维码登录")
const qrcodeStr = ref(null)
const showLoginWindow = ref(false)
const is_login = ref(null)

onMounted(() => {
    checkLogin();
}),

EventsOn("qrcodeStr", (qr)=>{
    console.log("收到qrcode");
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