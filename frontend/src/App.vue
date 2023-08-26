<template>
  <div class="main-container">
    <el-form :inline="true" :model="appForm" class="app-form" :rules="rules" ref="appFormRef">
      <el-form-item label="域名" prop="domain">
        <el-input class="form-item" v-model="appForm.domain" placeholder="请输入域名"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit(appFormRef)">探测</el-button>
      </el-form-item>
    </el-form>
    <div class="layout">
      <div class="basic-indicators">
        <h4>基本网络信息</h4>
        <p>
          DNS: <span>{{indicators.ip}}</span>
        </p>
        <div class="col-item">
          <p>Ping IP: <span>{{ indicators.ipDetectRet }}</span></p>
        </div>
        <div class="col-item">
          <p>TCP连接443端口: {{ indicators.tcpDetectRet }}</p>
        </div>
      </div>
      <div class="mtr">
        <h4>MTR结果列表</h4>
        <div v-if="skipText.length > 0">
          {{ skipText }}
        </div>
        <el-table
          v-else
          v-loading="mtrLoading"
          element-loading-text="请耐心等待，正尝试对ip进行mtr探测，可能需要几分钟时间"
          :data="mtrData">
          <el-table-column prop="target" label="Address" width="200"/>
          <el-table-column prop="Loss" label="Loss%">
            <template #default="scope">
              <span>{{ `${scope.row.loss_percent}%` }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="sent" label="Sent"/>
          <el-table-column label="Last(ms)">
            <template #default="{row}">
              {{ row.last_ms && row.last_ms.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="avg_ms" label="Avg(ms)">
            <template #default="{row}">
              {{ row.avg_ms ? row.avg_ms.toFixed(2) : 0 }}
            </template>
          </el-table-column>
          <el-table-column prop="best_ms" label="Best(ms)">
            <template #default="{row}">
              {{ row.best_ms ? row.best_ms.toFixed(2) : 0 }}
            </template>
          </el-table-column>
          <el-table-column prop="worst_ms" label="Worst(ms)">
            <template #default="{row}">
              {{ row.worst_ms ? row.worst_ms.toFixed(2) : 0 }}
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div> 
</template>

<script lang="ts" setup>
import {reactive, ref, onMounted} from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessageBox } from 'element-plus'
import { TryMTR } from '../wailsjs/go/detector/MTR'
import { StartTimer } from '../wailsjs/go/main/App'
import { TryPing } from '../wailsjs/go/detector/Ping'
import { ResolveHost } from '../wailsjs/go/detector/DNS'
import { TryHTTPSConnection } from '../wailsjs/go/detector/TCPConnection'

interface AppFormIn {
  domain: string
}

let mtrData = reactive<any[]>([])
const skipText = ref('')
const mtrLoading = ref(false)
const appForm = reactive<AppFormIn>({
  domain: ''
})
const appFormRef = ref<FormInstance>()
const rules = reactive<FormRules>({
  domain: [
    { required: true, message: '请输入域名', trigger: 'blur' }
  ],
})

const indicators = reactive({
  ip: '',
  ipDetectRet: '',
  tcpDetectRet: '',
  mtrDetectRet: '',
})

onMounted(() => {
  try {
    let appFormCache = localStorage.getItem('appForm')
    if (appFormCache) {
      const parsedAppFormCache: AppFormIn = JSON.parse(appFormCache) || {}
      if (parsedAppFormCache.domain) {
        appForm.domain = parsedAppFormCache.domain
        // windows开机启动时，网络连接可能存在一定延迟，导致探测失败，所以延迟2s
        setTimeout(async () => {
          await startDetect()
        }, 2000)
      }
    }
  } catch (e: any) {
    console.log('get appForm cache error', e)
  }
})

const onSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  const valid = await formEl.validate()
  if (!valid) return;

  localStorage.setItem('appForm', JSON.stringify(appForm))
  await startDetect()
}

const startDetect = async () => {
  // start timer
  StartTimer(appForm.domain)

  const ip = await ResolveHost(appForm.domain).catch((err) => {
    ElMessageBox.alert(err, '解析域名失败', {
      confirmButtonText: 'OK' as const
    })
  })

  if (ip) {
    indicators.ip = ip;
    const tasks = [detectIP(ip), detectTCP(ip), detectMTR(ip)]
    await Promise.all(tasks)
  }
}

const detectMTR = async (ip:string) => {
  mtrLoading.value = true

  try {
    const result = await TryMTR(ip)
    if (!result.HasPermission) {
      skipText.value = "权限不够跳过(需要管理员权限)"
      return
    }

    mtrData = result.Result

    console.log('Result:', result.Result)

    // check if result is Failed
    const hups = result.Result
    if ((hups.length === 0) || (hups[hups.length - 1].loss_percent >= 100)) {
      const err = `远程IP ${ip} 不可达, 请联系IT协助（可能是防火墙原因）`
      ElMessageBox.alert(err, 'MTR 失败', {
        confirmButtonText: 'OK' as const
      })
      return
    }
  } catch(err) {
    skipText.value = "skip"
    ElMessageBox.alert(err as string, 'MTR 失败', {
      confirmButtonText: 'OK' as const
    })
  } finally {
     mtrLoading.value = false
  }
}

const detectIP = async (ip:string) => {
  try {
    const result = await TryPing(ip)
    if (!result.HasPermission) {
      indicators.ipDetectRet = "权限不够跳过(需要管理员权限)"
      return
    }
    indicators.ipDetectRet = result?.Pass ? '正常' : '失败'
  } catch(err) {
    indicators.ipDetectRet = err as string
  }
}

const detectTCP = async (ip:string) => {
  try {
    await TryHTTPSConnection(ip)
    indicators.tcpDetectRet = '正常'
  } catch(err) {
    indicators.tcpDetectRet = `尝试HTTPS连接失败: ${err}`
  }
}
</script>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
.main-container {
  height: 100%;
  padding: 30px 15px;
}
.tip {
  margin-left: 40px;
  margin-bottom: 20px;
  font-size: 14px;
  color: #9d9d9d;
}
.layout {
  padding-bottom: 30px;
}
.form-item {
  width: 180px;
}
.el-form-item {
  margin-right: 15px!important;
}
.basic-indicators h4, .mtr h4 {
  margin-top: 8px;
  margin-bottom: 8px;
}
.basic-indicators p {
  margin-top: 5px;
  margin-bottom: 5px;
  font-size: 14px;
}
</style>
