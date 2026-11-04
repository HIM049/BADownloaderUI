import App from './App.vue'
import Varlet from '@varlet/ui'
import { createApp } from 'vue'
import '@varlet/ui/es/style'
import { GetTheme } from '../wailsjs/go/wails_api/WailsApi'
import { applyTheme } from './themes'

// 先应用配置中的主题再挂载，避免启动时颜色闪烁
GetTheme()
    .then(applyTheme, () => applyTheme())
    .finally(() => createApp(App).use(Varlet).mount('#app'))
