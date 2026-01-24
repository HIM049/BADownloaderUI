<template>
    <header class="pt-[10px]">
        <nav class="flex justify-between items-center mx-[13px]" style="--wails-draggable:drag;">
            <div id="logo" class="flex items-center">
                <div class="bg-white rounded-2xl p-2 flex">
                    <svg style="width: 40px;" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 396.65 396"><defs><clipPath id="a" transform="translate(-747 -329)"><rect x="850.03" y="452.03" width="174.85" height="176.95" style="fill:none" /></clipPath></defs><path d="M1089.72,430a8.3,8.3,0,0,1,8.28,8.28V642.72a8.3,8.3,0,0,1-8.28,8.28H785.28a8.3,8.3,0,0,1-8.28-8.28V438.28a8.3,8.3,0,0,1,8.28-8.28h304.44m0-30H785.28A38.28,38.28,0,0,0,747,438.28V642.72A38.28,38.28,0,0,0,785.28,681h304.44A38.28,38.28,0,0,0,1128,642.72V438.28A38.28,38.28,0,0,0,1089.72,400Z"transform="translate(-747 -329)" style="fill:#ffc7c7" /><line x1="128.5" y1="78.5" x2="114" y2="15" style="fill:none;stroke:#ffc7c7;stroke-linecap:round;stroke-miterlimit:10;stroke-width:30px" /><line x1="249.5" y1="78.5" x2="267" y2="15" style="fill:none;stroke:#ffc7c7;stroke-linecap:round;stroke-miterlimit:10;stroke-width:30px" /> <polyline points="331 213 330.37 361.58 331 386" style="fill:none;stroke:#333;stroke-linecap:round;stroke-linejoin:round;stroke-width:20px" /> <polyline points="386.65 319.24 331 386 270.29 319.24" style="fill:none;stroke:#333;stroke-linecap:round;stroke-linejoin:round;stroke-width:20px" /><g style="clip-path:url(#a)"><path d="M937.39,452.17c3.71,0,6.72,3.3,6.72,7.36V621.36c0,4.06-3,7.35-6.72,7.35s-6.72-3.29-6.72-7.35V459.53c0-4.06,3-7.36,6.72-7.36ZM897.05,481.6c3.72,0,6.73,3.29,6.73,7.35v103c0,4.07-3,7.36-6.73,7.36s-6.72-3.29-6.72-7.36V489c0-4.06,3-7.35,6.72-7.35Zm80.67,0c3.71,0,6.72,3.29,6.72,7.35v103c0,4.07-3,7.36-6.72,7.36S971,596,971,591.93V489c0-4.06,3-7.35,6.72-7.35ZM1018.05,511c3.72,0,6.73,3.29,6.73,7.35v44.14c0,4.06-3,7.36-6.73,7.36s-6.72-3.3-6.72-7.36V518.37c0-4.06,3-7.35,6.72-7.35Zm-161.33,0c3.71,0,6.72,3.29,6.72,7.35v44.14c0,4.06-3,7.36-6.72,7.36s-6.72-3.3-6.72-7.36V518.37c0-4.06,3-7.35,6.72-7.35Zm0,0" transform="translate(-747 -329)" style="fill:#8a8a8a" /></g> </svg>
                </div>
                <!-- Title Text -->
                <div class="flex text-xl font-extrabold ml-[10px]">
                    <p style="color: var(--color-primary-font);">B</p>
                    <p>ili</p>
                    <p style="color: var(--color-primary-font);">A</p>
                    <p>udio</p>
                    <p style="color: var(--color-primary-font);">D</p>
                    <p>ownloader</p>
                    <var-chip color="var(--color-body)" text-color="var(--color-text)" size="mini" class="ml-[10px]">{{ APP_VERSION }}</var-chip>
                </div>
            </div>
            <!-- Window Controls -->
            <div id="ctl-buttons" class="flex items-center" style="--wails-draggable:none;">
                <button class="bg-transparent border-none group cursor-pointer" id="minimize" @click="WindowMinimise()">
                    <svg class="h-[35px] w-[35px] p-[5px] ml-[5px] rounded group-hover:bg-[#acacac] group-hover:fill-white transition-colors duration-200" xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 -960 960 960"><path d="M200-440v-80h560v80H200Z" /></svg>
                </button>
                <button class="bg-transparent border-none group cursor-pointer" id="close" @click="Quit()">
                    <svg class="h-[35px] w-[35px] p-[5px] ml-[5px] rounded group-hover:bg-[#e94c4c] group-hover:fill-white transition-colors duration-200" xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 -960 960 960"> <path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z" /></svg>
                </button>
            </div>
        </nav>

        <slot />

    </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { WindowMinimise, Quit } from '../../../wailsjs/runtime'
import { GetAppVersion } from '../../../wailsjs/go/wails_api/WailsApi'

const APP_VERSION = ref(null)
onMounted(() => {
    GetAppVersion().then(result => {
        APP_VERSION.value = result
    })
})
</script>
