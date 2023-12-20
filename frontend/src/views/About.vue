<template #title>
  <el-button style="color:#1f8dd6" text @click="dialog = true" >模板列表</el-button>

  <el-drawer
      :with-header="false"
      :visible.sync="drawer.config"
      :title="title"
      :before-close="handleClose"
      v-model="dialog"
      size="80%"

      class="page-cantor"
  >
  </el-drawer>
  <el-form :model="form" label-width="-10px">

    <el-form-item >
<!--      <el-input class="DocSearch-Input" placeholder="Search" v-model="form.query" @keyup.enter="eventEnter"/>-->
<!--      <el-button type="primary" @click="begin">Search </el-button>-->
      <n-input v-model="form.query"  placeholder="请输入关键字" @keyup.enter="eventEnter"/>
      <n-button @click="begin">搜索</n-button>
    </el-form-item>



    <el-button-group style="display: flex;justify-content: center;align-items: center;" v-model="page"  >
      <el-button v-show="dialogVisible1 === true"  @click="prePage" >上一页</el-button>
      <el-button v-show="dialogVisible2  === true"   @click="nextPage "  >下一页<el-icon class="el-icon--right"><ArrowRight /></el-icon></el-button>
    </el-button-group>


    <div align="center"><h1>共为你找到相关结果约为{{ hits }}个。显示从{{ start }}起共{{ listLength }}个。</h1></div>

    <div class="table-responsive-vertical shadow-z-1">
      <table id="table" class="table table-striped table-hover table-mc-indigo">
        <tbody>
        <div id="video-container">
          <div class="video-item" v-for="item in list" :key="item.Id"  >
            <div class="video-val">
              <video class="video-player item.Payload.Number" controls :poster="item.Payload.Pic" preload="metadata" @mouseenter="playVideo()" @mouseleave="stopVideo()">
                <source class="video-source" :src="item.Payload.VideoShort" type="video/mp4">
              </video>

            </div>
            <div class="video-title">
              <el-link class="el-link" type="primary" :href="item.Payload.VideoPreview" target="_blank">{{ item.Payload.Title }}</el-link>
              <el-link class="el-label" :href="item.Payload.Url" target="_blank"> 视频编号 - {{item.Payload.Number}}</el-link>
              <el-button type="primary" @click="openFolder(item.Payload.FormalLocal)">Open Folder</el-button>
            </div>
          </div>
        </div>


        </tbody>
      </table>
    </div>

    <el-button-group style="display: flex;justify-content: center;align-items: center;" v-model="page"  >
      <el-button v-show="dialogVisible1 === true"  @click="prePage"  >上一页</el-button>
      <el-button v-show="dialogVisible2  === true"   @click="nextPage "  >下一页<el-icon class="el-icon--right"><ArrowRight /></el-icon></el-button>
    </el-button-group>

  </el-form>

</template>


<script>
// import { ExportVideoList, VideoList, VideoManage } from '../../wailsjs/go/main/App'
import { VideoList,SearchKey,OpenFolder } from "../../wailsjs/go/internal/App.js";

import { defineComponent, reactive, ref, onMounted, computed, h } from 'vue'
import { ElNotification, ElDrawer, ElMessageBox } from "element-plus"

export default defineComponent({
  setup() {

    const dialogVisible1 = ref(true)
    const dialogVisible2 = ref(true)

    const datas = reactive({
      config: {
        username: "admin",
        password: "123456",
      },
    })
    const drawer = reactive({
      config: false,
    })
    const dialog = ref(false)

    const form = reactive({})
    let list = ref()
    let query = ref()
    let hits = ref()
    let start = ref()
    let prevFrom = ref()
    let nextFrom = ref()
    let listLength = ref()
    let page = ref()

    const showButton = () => {
      if (prevFrom.value < 0) {
        dialogVisible1.value = false
      } else {
        dialogVisible1.value = true
      }

      if (nextFrom.value === hits.value) {
        dialogVisible2.value = false
      } else {
        dialogVisible2.value = true
      }
    }

    const prePage = () => {
      form.page = prevFrom.value
      searchVideoList()
    }

    const nextPage = () => {
      form.page = nextFrom.value
      searchVideoList()
    }

    const begin = () => {

      hits.value = 0
      start.value = 0
      form.page = 0
      // searchVideoList()
      // alert(form)
      searchKey(form)
    }

    const searchVideoList = async () => {
      // alert("执行了 searchVideoList")
      if (form.query === undefined) {
        form.query = query.value
      }
      await searchKey()
    }

    const searchKey = async (form) =>{
      SearchKey(form).then(res => {
        if (res.code !== 200) {
          ElNotification({
            title: res.msg,
            type: "error",
          })
          return
        }
        list.value = res.data.Items
        listLength.value = list.value.length;
        query.value = res.data.Query
        hits.value = res.data.Hits
        start.value = res.data.Start
        prevFrom.value = res.data.PrevFrom
        nextFrom.value = res.data.NextFrom

        showButton()
      })
    }

    const openFolder = (path) => {
      OpenFolder(path).then(res => {
        if (res.code !== 200) {
          ElNotification({
            title: res.msg,
            type: "error",
          })
        }
      })
    }

    const videoList = () => {
      VideoList().then(res => {
        if (res.code !== 200) {
          ElNotification({
            title: res.msg,
            type: "error",
          })
        }
        list.value = res.data.Items
        listLength.value = list.value.length;
        query.value = res.data.Query
        hits.value = res.data.Hits
        start.value = res.data.Start
        prevFrom.value = res.data.PrevFrom
        nextFrom.value = res.data.NextFrom

        showButton()
      })
    }

    onMounted(() => {
      searchVideoList()
    })

    return {
      list,
      page,
      drawer,
      form,
      dialog,
      hits,
      start,
      prevFrom,
      nextFrom,
      listLength,
      dialogVisible1,
      dialogVisible2,
      showButton,
      prePage,
      nextPage,
      begin,
      searchVideoList,
      openFolder,
      videoList,
      searchKey,
      handleClose() {
        // 处理关闭事件
      },
      eventEnter(){
        alert(1)
      },
      playVideo() {
        const video = event.target;
        const source = video.querySelector(".video-source");

        if (!source.dataset.src) {
          source.dataset.src = source.src; // 使用自定义属性保存视频的源路径
          video.load(); // 加载视频
        }
        video.playbackRate = 1;
        // video.currentTime = 10; // 初识开始秒数
        video.muted = false; // true 静音  false 播放
        video.play();
      },

      stopVideo() {
        const video = event.target;
        video.pause(); // 暂停视频播放
        video.currentTime = 0; // 重置视频播放时间
      },
    }
  }
})


</script>



<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

/* 设置视频列表样式 */
.video-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-bottom: 20px;
}

/* 设置视频缩略图样式 */
.video-player {
  width: 200px;
  height: 150px;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
  margin-right: 20px;
}

/* 设置视频标题样式 */
.video-title {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

/* 设置视频标题链接样式 */
.el-link {
  font-size: 18px;
  color: #0078d7;
  margin-bottom: 10px;
}
.el-label {
  font-size: 18px;
  color: #0078d7;
  margin-bottom: 10px;
}


/* 设置分页信息样式 */
h1 {
  font-size: 18px;
  color: #0078d7;
  margin-bottom: 20px;
}
.DocSearch-Input {
  appearance: none;
  background: transparent;
  border: 0;
  flex: 1;
  font: inherit;
  font-size: 1.2em;
  height: 100%;
  outline: none;
  padding: 0 0 0 8px;
  width: 50%;
  margin-top: 10px;
}
</style>
