<template #title>
  <el-button style="color:#1f8dd6" text @click="dialog = true" >用户信息</el-button>

  <el-drawer
        :with-header="false"
        :visible.sync="drawer.config"
        :title="title"
        :before-close="handleClose"
        v-model="dialog"
        size="39%"
        class="config-area"
    >
      <el-form ref="form" :model="datas.config"  label-width="60px">

        <el-form-item label="账号" prop="username">
          <el-input v-model="datas.config.username" placeholder="admin"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="ps">
          <el-input v-model="datas.config.password" type="password" placeholder="123456"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="goUpdateConfig" :loading="loading"> 配 置 </el-button>
          <el-button type="danger" @click="dialog = false"> 取 消 </el-button>
        </el-form-item>
      </el-form>
    </el-drawer>


</template>

<script lang="ts" setup >
import {ref, reactive, watch, watchEffect} from "vue";

import {ConfigEdit} from "../../wailsjs/go/internal/App";

import {ElDrawer, ElMessageBox, ElNotification} from "element-plus";
// let props = defineProps(['title', 'btnType', 'data'])

const datas = reactive({
  config: {
    username :"admin",
    password :"123456",
  },
})
const drawer = reactive({
  config: false,
})
const dialog = ref(true)
const loading = ref(false)
let timer
const formLabelWidth = '80px'

const cancelForm = () => {
  loading.value = false
  dialog.value = false
  clearTimeout(timer)
}
const rules = reactive({
  username: [
    { required: true, message: 'Please input Activity name', trigger: 'blur' },
    { min: 3, max: 5, message: 'Length should be 3 to 5', trigger: 'blur' },
  ],})

const drawerRef = ref<InstanceType<typeof ElDrawer>>()
const onClick = () => {
  drawerRef.value!.close()
}


const handleClose = (done) => {
  if (loading.value) {
    return
  }

  // ElMessageBox.confirm("").then(() => {
  //       loading.value = true
  //       timer = setTimeout(() => {
  //         done()
  //         // 动画关闭需要一定的时间
  //         setTimeout(() => {
  //           loading.value = false
  //         }, 0)
  //       }, 0)
  //     })
  //     .catch(() => {
  //       // catch error
  //     })
}


function goUpdateConfig() {
  // dialog.value = false

  ConfigEdit(datas).then(res => {
    if (res.code !== 200) {ElNotification({
      title:res.msg,
      type: "error",
    })
      dialog.value = true
      return
    }
    dialog.value = false
    loading.value = false
    ElNotification({
      title:res.msg,
      type: res.data,
    })

    datas.config.username = res.data.UserName
    datas.config.password = res.data.PassWord
  })

}

</script>


<style scoped>

</style>