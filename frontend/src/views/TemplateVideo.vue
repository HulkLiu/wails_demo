<template>
  <el-container class="page-cantor" >
    <el-header class="header-area" style="height: 20px;">

    </el-header>
    <el-aside class="side-area" style="width: 78px;--wails-draggable: drag;">
    </el-aside>

    <el-main class="main-area" style="--wails-draggable:no-drag">

      <el-template>
        <n-form>
          <n-input v-model:value="query"  placeholder="请输入关键字" @keyup.enter="search"/>
          <n-button @click="search">搜索</n-button>

          <el-button-group style="display: flex;justify-content: center;align-items: center;" v-model:value="page"   >
            <el-button v-show="dialogVisible1 === true"  @click="prePage" >上一页</el-button>
            <el-button v-show="dialogVisible2  === true"   @click="nextPage "  >下一页</el-button>
          </el-button-group>

          <div align="center"><h1>共为你找到相关结果约为{{ hits }}个。显示从{{ start }}起共{{ listLength }}个。</h1></div>

          <div class="table-responsive-vertical shadow-z-1">
            <table id="table" class="table table-striped table-hover table-mc-indigo">
              <tbody>
              <div id="video-container" >
                <div class="video-item"  v-for="item in list" :key="item.Id"  >
                  <div class="video-val">
                    <video class="video-player item.Payload.Number" controls :poster="item.Payload.Pic" preload="metadata" @mouseenter="playVideo()" @mouseleave="stopVideo()">
                      <source class="video-source" :src="item.Payload.VideoShort" type="video/mp4">
                    </video>
                  </div>

                  <div class="video-title">
                    <el-link class="el-link" title="视频小样链接" type="primary" :href="item.Payload.VideoPreview" target="_blank">{{ item.Payload.Title }}</el-link>
                    <el-link class="el-label" title="网站链接" :href="item.Payload.Url" target="_blank"> 视频编号 - {{item.Payload.Number}}</el-link>
                    <el-button type="primary" @click="openFolder(item.Payload.FormalLocal)">打开文件</el-button>
                  </div>


                </div>
              </div>


              </tbody>
            </table>
          </div>

          <el-button-group style="display: flex;justify-content: center;align-items: center;" v-model:value="page"   >
            <el-button v-show="dialogVisible1 === true"  @click="prePage" >上一页</el-button>
            <el-button v-show="dialogVisible2  === true"   @click="nextPage "  >下一页</el-button>
          </el-button-group>
        </n-form>
      </el-template>
    </el-main>

  </el-container>

</template>

<script>
import {defineComponent, onMounted, ref} from "vue";
import {OpenFolder, SearchKey} from "../../wailsjs/go/internal/App.js";
import {ElNotification} from "element-plus";

export default defineComponent({
  setup () {
    const dialogVisible1 = ref(true)
    const dialogVisible2 = ref(true)
    let form = ref({
      page:"",
      query:"",
    })
    let list = ref()
    let page = ref()
    let query = ref()
    let hits = ref()
    let start = ref()
    let prevFrom = ref()
    let nextFrom = ref()
    let listLength = ref()

    onMounted(() => {
      searchKey()
    })

    const showButton = () => {
      dialogVisible1.value = prevFrom.value >= 0;
      dialogVisible2.value = nextFrom.value !== hits.value;
    }

    const prePage = () => {
      form.query = query.value
      form.page = prevFrom.value
      searchKey(form)
    }

    const nextPage = () => {
      form.query = query.value
      form.page = nextFrom.value
      searchKey(form)
    }
    const searchVideoList = async () => {
      // alert("执行了 searchVideoList")
      if (form.query === undefined) {
        form.query = query.value
      }
      await searchKey(form)
    }

     const searchKey =  (form) =>{
      console.log(form)

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
        // form.query = res.data.Query


        showButton()
      })
    }

    return {

      form,
      list,
      page,
      listLength,
      hits,
      start,
      query,
      prevFrom ,
      nextFrom ,
      dialogVisible1,
      dialogVisible2,
      search(){
        // alert(query.value)
        form.query= query.value
        form.page = page.value
        // console.log(form)

        SearchKey(form).then(res => {
          console.log(form)

          if (res.code !== 200) {
            ElNotification({
              title: res.msg,
              type: "error",
            })
            return
          }
          list.value = res.data.Items
          listLength.value = list.value.length;
          hits.value = res.data.Hits
          start.value = res.data.Start
          prevFrom.value = res.data.PrevFrom
          nextFrom.value = res.data.NextFrom
          query.value = res.data.Query

          showButton()
        })

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

      openFolder(path) {
        OpenFolder(path).then(res => {
          if (res.code !== 200) {
            ElNotification({
              title: res.msg,
              type: "error",
            })
          }
        })
      },
      prePage(){
        prePage()

      },
      nextPage(){
        nextPage()
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
.video-val {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
/* 设置视频标题样式 */
.video-title {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
/* 设置视频标题样式 */
.video-description {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  width: 60%;
  font-size: 18px;
  color: #0078d7;
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
