<template>
    <n-button-group>

      <n-button secondary strong :render-icon="renderIcon" @click="activate('right')">账户设置</n-button>

    </n-button-group>

    <n-drawer v-model:show="active" :width="502" :placement="placement" :on-after-leave="leaveDrawer">
      <n-drawer-content title="信息登录">
        <n-form
            ref="formRef"
            :model="datas"
            label-placement="left"
            :label-width="89"
            :style="{maxWidth: '500px'}"
            :rules="rules"
        >
          <n-form-item label="部门" path="config.repo">
            <n-input  v-model:value="datas.config.repo" placeholder="请输入部门"></n-input>
          </n-form-item>
          <n-form-item label="邮箱/电话" path="config.email">
            <n-input  v-model:value="datas.config.email" placeholder="请输入邮箱/电话"></n-input>
          </n-form-item>
          <n-form-item label="账号" path="config.owner">
            <n-input  v-model:value="datas.config.owner" placeholder="请输入账号"></n-input>
          </n-form-item>

          <n-form-item label="密码" path="config.access_token">
            <n-input type="password"  v-model:value="datas.config.access_token"></n-input>
          </n-form-item>
          <n-form-item>
            <n-row :gutter="[0, 24]">
              <n-col :span="24">
                <div style="display: flex; justify-content: flex-end">
                  <n-button type="primary" @click="goUpdateConfig"
                            :disabled="datas.config.repo === '' ||
                            datas.config.email === '' ||
                            datas.config.owner === '' ||
                            datas.config.access_token === ''"
                            :loading="loading">
                    设置
                  </n-button>
                  <n-button type="danger" @click="drawer.config = false">
                    取消
                  </n-button>
                </div>
              </n-col>
            </n-row>

          </n-form-item>
        </n-form>

      </n-drawer-content>
    </n-drawer>

</template>

<script>
import { CashOutline as CashIcon } from "@vicons/ionicons5";
import { NIcon } from "naive-ui";
import {computed, defineComponent, onMounted, ref, h} from "vue";
import {GetConfig, SetConfig} from "../../wailsjs/go/internal/App.js";
import {ElNotification} from "element-plus";
import {useMessage} from "naive-ui";

export default defineComponent({
  components: {
    CashIcon
  },
  setup() {
    const formRef = ref(null)
    const drawer = ref({
      config:true
    })
    const datas = ref({
      config: {
        repo:"",
        email:"",
        owner:"",
        access_token:"",
      },
    })

    const goGetConfig = () =>{
      GetConfig().then(res => {
        // console.log(res.data)
        datas.value.config = res.data;
        if (datas.value.config.repo === '') {
          active.value = true
        }else{
          active.value = false
        }
      })
    }
    const setConfig = () =>{
      if (datas.value.config.owner === '') {
        return
      }
      SetConfig( JSON.stringify(datas.value.config) ).then(res => {
        if (res.code !== 0) {
          message.error(res.msg)
          return;
        }
        message.success("设置成功")
        // drawer.value.config = false;
        active.value = false
      })
    }

    onMounted(() => {
      activate('right')
      goGetConfig()
    })
    const message = useMessage()
    const active = ref(true);
    const placement = ref("right");
    const activate = (place) => {
      active.value = true;
      placement.value = place;
    };

    return {
      formRef,
      active,
      placement,
      activate,
      drawer,
      datas,
      goGetConfig,
      goUpdateConfig(){

        setConfig()
      },
      leaveDrawer(){
        if ( datas.value.config.repo === '' || datas.value.config.owner === '' || datas.value.config.email === '' || datas.value.config.access_token === '') {
          message.warning("请配置登录信息！")
          active.value = true
        }
      },
      rules: {
        config:{
            repo: {
              required: true,
              message: '请输入部门',
              trigger: 'blur'
            },
            email: {
              required: true,
              message: '请输入邮箱/电话',
              trigger: 'blur'
            },
            owner: {
              required: true,
              message: '请输入账号',
              trigger: 'blur'
            },
            access_token: {
              required: true,
              message: '请输入密码',
              trigger: 'blur'
            },


        },

      },
      renderIcon() {
        return h(NIcon, null, {
          default: () => h(CashIcon)
        });
      },

    };
  }
});

</script>

<style scoped>

</style>