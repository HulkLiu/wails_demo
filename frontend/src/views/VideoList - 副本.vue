<template>
    <el-container class="page-cantor" v-loading="loading" :style="areaStyle">
      <el-header class="header-area" style="height: 30px;">
        <div>
          筛选框位置

          <n-button @click="VideoListExport">导出 EXCEL </n-button>
          <n-modal :show="showModal">
            <n-card
                style="width: 600px"
                title="模态框"
                :bordered="false"
                size="huge"
                role="dialog"
                aria-modal="true"
            >
              <template #header-extra>
                <video width="80" height="80" controls  id="videoPlayer">
                  <source src="" type="video/mp4" id="videoSrc">
                </video>
              </template>

            </n-card>
          </n-modal>

        </div>
      </el-header>
      <el-aside class="side-area" style="width: 78px;--wails-draggable: drag;">
      </el-aside>
      <el-main class="main-area" style="--wails-draggable:no-drag">
        <el-template>
          <n-input v-model:value="keyword" placeholder="请输入关键字" @keyup.enter="search"/>
          <n-button @click="search">搜索</n-button>
          <n-p> 你选中了 {{ checkedRowKeys.length }} 行。 </n-p>
          <n-data-table
              :key="(row) => row.key"
              :columns="columns"
              :data="data"
              :row-key="rowKey"
              :pagination="paginationRef"
              :on-update:page="handlePageChange"
              @update:checked-row-keys="handleCheck"
          />
        </el-template>
      </el-main>

    </el-container>

</template>

<script>
import {h, defineComponent, ref, nextTick, computed, onMounted} from 'vue'
import { NInput,NTag, NButton, useMessage , useDialog } from 'naive-ui'
import {VideoManage,ExportVideoList} from "../../wailsjs/go/internal/App.js";
import {ElNotification} from "element-plus";


const ShowOrEdit = defineComponent({
  props: {
    value: [String, Number],
    onUpdateValue: [Function, Array]
  },
  setup (props) {
    const isEdit = ref(false)
    const inputRef = ref(null)
    const inputValue = ref(props.value)
    function handleOnClick () {
      isEdit.value = true
      nextTick(() => {
        inputRef.value.focus()
      })
    }
    function handleChange () {
      props.onUpdateValue(inputValue.value)
      isEdit.value = false
    }

    return () =>
        h(
            'div',
            {
              style: 'min-height: 22px',
              onClick: handleOnClick
            },
            isEdit.value
                ? h(NInput, {
                  ref: inputRef,
                  value: inputValue.value,
                  onUpdateValue: (v) => {
                    inputValue.value = v
                  },
                  onChange: handleChange,
                  onBlur: handleChange
                })
                : props.value
        )
  }
})
// let showModal= ref(true)

export default defineComponent({
  setup () {
    // const data = ref(createData())
    const data = ref()
    const page = ref(1)
    const keyword = ref()
    // const dialog1 = useDialog()
    const getDataIndex = (key) => {
      return data.value.findIndex((item) => item.key === key)
    }
    const handlePageChange = (curPage) => {
      page.value = curPage
    }

    const paginationRef = computed(() => ({
      pageSize: 10,
      page: page.value
    }))

    const searchData = () => {
        VideoManage().then(res => {
          if (res.code !== 200) {
            ElNotification({
              title:res.msg,
              type:"error",
            })
          }
          data.value = res.data.Items
          console.log(data.value )
          // createData(res.data.Items)
        })
    }

    const checkedRowKeysRef = ref([]);
    onMounted(() => {
      searchData()
    })

    return {
      showModal: ref(false),
      keyword,
      searchData,
      data,
      paginationRef,
      checkedRowKeys: checkedRowKeysRef,
      handlePageChange,
      rowKey: (row) => row,
      handleCheck(rowKeys) {
        // alert(rowKeys[0].age)
        console.log(rowKeys)
        checkedRowKeysRef.value = rowKeys;
      },
      ControlVideo(src) {
        alert(src)
        showModal = true
        // alert(showModal)
        PlayShortVideo(src)
      },
      VideoListExport(){
        ExportVideoList().then(res => {
          if (res.code !== 200) {
            ElNotification({
              title:res.msg,
              type:"error",
            })
          }
          // alert(res.data)
          // createData(res.data.Items)
        })
      },

      search(){
        VideoManage().then(res => {
          if (res.code !== 200) {
            ElNotification({
              title:res.msg,
              type:"error",
            })
          }
          data.value = res.data.Items
          console.log(data.value )
          // createData(res.data.Items)
        })
      },
      columns: [
        {
          type: "selection",
          disabled(row) {
            return row.name === "Edward King 3";
          }
        },
        // {
        //   title:"图片",
        //   key: 'Pic',
        //
        //   render (row) {
        //     return h(
        //         NButton,
        //         {
        //           size: 'small',
        //           onClick: () => ControlVideo(row.Payload.VideoShort ) ,
        //
        //         },
        //         { default: () => '播放视频' }
        //     )
        //   }
        // },
        // {
        //   type:"video",
        //   title:"视频",
        //
        // },
        // {
        //   title: 'Number',
        //   key: 'Number',
        //   width: 100,
        //   render (row) {
        //     const index = getDataIndex(row.key)
        //     return h(ShowOrEdit, {
        //       value: row.Payload.Number,
        //
        //       onUpdateValue (v) {
        //         data.value[index].Payload.number = v
        //       }
        //     })
        //   }
        // },
        {
          title: '标题',
          key: 'Title',
          width: 200,
          render (row) {
            const index = getDataIndex(row.key)
            return h(ShowOrEdit, {
              value: row.Payload.Title,
              onUpdateValue (v) {
                data.value[index].Payload.Title = v
              }
            })
          }
        },
        {
          title: '网站链接',
          key: 'Url',
          width: 200,
          render (row) {
            const index = getDataIndex(row.key)
            return h(ShowOrEdit, {
              value: row.Payload.Url,
              onUpdateValue (v) {
                data.value[index].Payload.Url = v
              }
            })
          }
        },
        {
          title: '本地链接',
          key: 'ShortLocal',
          render (row) {
            const index = getDataIndex(row.key)
            return h(ShowOrEdit, {
              value: row.Payload.ShortLocal,
              onUpdateValue (v) {
                data.value[index].Payload.ShortLocal = v
              }
            })
          }
        }
      ]
    }
  }
})
async function PlayShortVideo(src) {

  document.getElementById('videoSrc').src = src;
  let player = document.getElementById('videoPlayer');
  await player.load();
  player.play();
}

</script>