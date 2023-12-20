<script>
import {computed, h, ref} from "vue";
import {RouterLink, useRouter, useRoute} from "vue-router";
import { NIcon, NConfigProvider, NMenu, NTabs, NTabPane, NSpace, NLayout, NLayoutSider, NLayoutContent, NMessageProvider, NSwitch, NDivider } from 'naive-ui';
// import { NIcon, NConfigProvider } from "naive-ui";
import {
  Home as HomeIcon,
  AlertCircleSharp as AboutIcon,
  SettingsSharp as SettingIcon,
  Barcode as OcrIcon,


} from '@vicons/ionicons5'

import { darkTheme } from 'naive-ui'
import mitt from './utils/event.js'
import Login from "./views/Login.vue";

const renderIcon = (icon) => {
  return () => h(NIcon, null, { default: () => h(icon) });
}



const menuOptions =  [

  {
    label: () => h(
        RouterLink,
        {
          to: {
            path: "/",
          }
        },
        { default: () => "主页" }
    ),
    key: "主页",
    icon: renderIcon(HomeIcon),
    path:"/"
  },
  {
    label: () => h(
        RouterLink,
        {
          to: {
            name: 'about',
            path: "/about",
          }
        },
        { default: () => "模板预览" }
    ),
    key: "模板预览",
    icon: renderIcon(AboutIcon),
    path: "/about",

  },
  {
    label: () => h(
        RouterLink,
        {
          to: {
            name: 'videoList',
            path: "/videoList",
          }
        },
        { default: () => "视频管理" }
    ),
    key: "视频管理",
    icon: renderIcon(OcrIcon),
    path: "/videoList",
  },
  {
    label: () => h(
        RouterLink,
        {
          to: {
            name: 'task',
            path: "/task",
          }
        },
        { default: () => "任务管理" }
    ),
    key: "任务管理",
    icon: renderIcon(OcrIcon),
    path: "/task",
  },
  {
    label: () => h(
        RouterLink,
        {
          to: {
            name: 'setting',
            path: "/setting",
          }
        },
        { default: () => "软件设置" }
    ),
    key: "软件设置",
    icon: renderIcon(SettingIcon),
    path: "/setting",

  },
  // {
  //   label: () => h(
  //       RouterLink,
  //       {
  //         to: {
  //           name: 'test',
  //           path: "/test",
  //         }
  //       },
  //       { default: () => "测试" }
  //   ),
  //   key: "测试",
  //   icon: renderIcon(OcrIcon),
  //   path: "/test",
  // },

]


export default {
  components: {
    Login,
    NConfigProvider
  },

  data() {
    const valueRef = ref();
    const panelsRef = ref([]);
    const addableRef = computed(() => {
      return {
        disabled: panelsRef.value.length >= 10
      };
    });
    const closableRef = computed(() => {
      return panelsRef.value.length > 1;
    });
    const router = useRouter();
    // const router = useRoute();
    const activeTab = ref(); // Use the current route path as the active tab

    const switchTheme = ref(false);
    const myTheme = ref(null);

    function changeTheme() {
      myTheme.value = switchTheme.value ? darkTheme : null;
    }

    const handleTabChange = (key) => {

      let p = ""
      menuOptions.forEach(option => {
        if ( key === option.key ) {
          p = option.path
        }
      });

      console.log(p)
      activeTab.value = key;
      router.push(p);

    };
    const panels = ref([]);
    const tabPanes = ref([]);

    const handleTabClose = (key) => {
      // console.log(key)
      const index = panels.value.findIndex(tabPane => tabPane.name === key);
      if (index !== -1) {
        panels.value.splice(index, 1);
        if (activeTab.value === key) {
          if (panels.value.length) {
            activeTab.value = panels.value[Math.max(0, index - 1)].name;
            handleTabChange(activeTab.value)
          } else {
            activeTab.value = '';
          }
        }
      }
    };
    return {
      handleTabClose,
      handleTabChange,
      menuOptions :menuOptions,
      dialogVisible : true,
      collapsed: false,
      switchTheme: false,
      myTheme: null,
      rt: window.runtime,
      activeTab,
      value: valueRef,
      panels,
      tabPanes,
      // panels: panelsRef,
      addable: addableRef,

      changeTheme,
      railStyle: ({
        focused,
        checked
      }) => {
        const style = {};
        if (checked) {
          style.background = "#4B9D5F";
          if (focused) {
            style.boxShadow = "0 0 0 2px #d0305040";
          }
        } else {
          style.background = "#000000";
          if (focused) {
            style.boxShadow = "0 0 0 2px #2080f040";
          }
        }
        return style;
      },

      handleUpdateValue(key, item) {

        if (!panels.value.find(tabPane => tabPane.name === item.key)) {
          panels.value.push({
            label: key.charAt(0).toUpperCase() + key.slice(1),
            name: key,

          });
        }
        console.log(panels.value)
        valueRef.value = key
        activeTab.value = key

      },
    }
  },
  mounted () {
    // 初始化项目时将主题保存在localStorage中
    localStorage.setItem('theme', 1)
  },

  methods: {
    changeTheme() {
      if (this.switchTheme) {
        this.myTheme = darkTheme
        localStorage.setItem("theme", 0)
        mitt.emit("theme","0")
      } else {
        this.myTheme = null
        localStorage.setItem("theme", 1)
        mitt.emit("theme","1")
      }
    }
  }
}

</script>

<template>
  <n-config-provider :theme="myTheme">
    <n-space vertical size="large">
      <n-layout has-sider position="absolute">

        <n-layout-sider bordered collapse-mode="width" :collapsed-width="80" :width="150" :collapsed="collapsed"
          show-trigger @collapse="collapsed = true" @expand="collapsed = false" style="--wails-draggable:drag; opacity: 1;">
            <div align="center" style="margin-top: 10px">
              <n-message-provider>
                <Login/>
              </n-message-provider>
            </div>
            <n-menu  v-model:value="activeTab"  @update:value="handleUpdateValue"  :options="menuOptions" :collapsed-width="64" :collapsed-icon-size="22" style="margin-top: 40px;" />
            <div class="switchBtnPar">
              <n-divider />
              <n-switch :rail-style="railStyle" v-model:value="switchTheme" @update:value="changeTheme()"
                class="switchBtn">
                <template #checked>
                  亮
                </template>
                <template #unchecked>
                  暗
                </template>
              </n-switch>
            </div>


        </n-layout-sider>
        <n-layout-content>
          <n-message-provider>
            <n-tabs
                tab-style="min-width: 80px;"
                type="card"
                v-model:value="activeTab"
                @update:value="handleTabChange"
                closable
                @close="handleTabClose"
            >

              <n-tab-pane v-for="item in panels" :key="item.name" :name="item.name" :label="item.label">

              </n-tab-pane>

            </n-tabs>

            <keep-alive>
              <router-view />
            </keep-alive>

          </n-message-provider>

        </n-layout-content>
      </n-layout>
    </n-space>
  </n-config-provider>
</template>

<style>
.switchBtnPar {
  position: relative;
}

.switchBtn {
  position: absolute;
  left: 50%;
  transform: translate(-50%);
}
body {
  margin: 0;
}
</style>