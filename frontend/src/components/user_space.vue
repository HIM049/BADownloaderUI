<template>
    <FramePage title="个人空间">
        <h3>登录</h3>
        {{ loginText }}
        <img :src="qrcodeStr" alt="Image">
        <var-space justify="flex-end">
            <var-button type="primary" @click="login">登录</var-button>
        </var-space>
    </FramePage>
</template>

<script setup>
import FramePage from '../components/modules/frame_page.vue'
import { ref } from 'vue'
import { LoginBilibili } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOnce } from '../../wailsjs/runtime'

const loginText = ref("请扫描二维码登录")
const qrcodeStr = ref(null)

EventsOn("qrcodeStr", (qr)=>{
    console.log("收到qrcode");
    qrcodeStr.value = "data:image/png;base64," + qr
})

EventsOn("loginStatus", (status) => {
    console.log(status);
    loginText.value = status
})

function login() {
    LoginBilibili()
}
</script>