<template>
  <el-container class="page-cantor" v-loading="loading" :style="areaStyle">
    <el-header class="header-area" style="height: 89px;">
      <n-button @click="VideoListExport">导出 EXCEL </n-button>
      <n-button @click="dialogVisibleAdd"> 新增视频 </n-button>
      <n-input v-model:value="keyword" placeholder="请输入关键字" @keyup.enter="search" @emit-video-list="search"/>
    </el-header>
    <n-message-provider>

      <el-main class="main-area" style="--wails-draggable:no-drag">

      <n-data-table
          :columns="cols"
          :data="data"
          :row-props="rowProps"
          :pagination="paginationRef"
      />
      <n-dropdown
          placement="bottom-start"
          trigger="manual"
          :row-key = "rowKey"
          :x="x"
          :y="y"
          :options="options"
          :show="showDropdown"
          :on-clickoutside="onClickoutside"
          @select="handleSelect"
      />

      <n-modal v-model:show="dialogVisible" style="width: 1500px">
        <n-card
            style="width: 67%"
            title="视频新增/修改"
            :bordered="false"
            size="huge"
            role="dialog"
            aria-modal="true"

            :style="{maxWidth: '840px'}"
            :label-width="90"
        >

          <template #header-extra></template>

            <n-form  :model="form" :rules="rules"  label-placement="left">
              <n-alert title="新增文件规则" type="info" :bordered="false" style="margin-bottom :30px;">
                1、新增视频前：请创建如下目录 <br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Z:\华章素材库 => 视频编号名称文件夹 => 4种资源文件<br><br>
                2、视频编号规则参考：https://www.vjshi.com/watch/27067568.html <br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;视频编号：27067568
              </n-alert>
              <n-form-item v-show="false" path="Type" label="Type">
                <n-input :disabled="true" v-model:value="form.Type" placeholder="自动添加" />
              </n-form-item>
              <n-form-item v-show="false" path="Id" label="Id" >
                <n-input :disabled="true" v-model:value="form.Id" placeholder="自动添加" />
              </n-form-item>

              <n-form-item  path="Number" label="视频编号 " title="aa">
                <n-button @click="autoGetVideoNumber">自动获取</n-button>
                <n-input v-model:value="form.Payload.Number" placeholder="新增时，请保持编号视频不要重复！" />
              </n-form-item>
              <n-form-item path="Title" label="视频标题">
                <n-input v-model:value="form.Payload.Title" placeholder="视频标题" />
              </n-form-item>
              <n-form-item path="ShortLocal" label="本地预览地址">
                <n-input v-model:value="form.Payload.ShortLocal" placeholder="G:\Videos/8909441/滴水穿石 岩石滴水_8909441_VideoShort.mp4" />
              </n-form-item>
              <n-form-item path="FormalLocal" label="本地资源地址">
                <n-input v-model:value="form.Payload.FormalLocal" placeholder="G:\Videos/8909441/滴水穿石 岩石滴水_8909441_VideoFormal.mp4" />
              </n-form-item>
              <n-form-item path="PreviewLocal" label="本地小样地址">
                <n-input v-model:value="form.Payload.PreviewLocal" placeholder="G:\Videos/8909441/滴水穿石 岩石滴水_8909441_VideoPreview.mp4" />
              </n-form-item>
              <n-form-item path="ImgLocal" label="本地图片地址">
                <n-input v-model:value="form.Payload.ImgLocal" placeholder="G:\Videos/8909441/滴水穿石 岩石滴水_8909441.jpg" />
              </n-form-item>

              <n-form-item path="Url" label="视频地址">
                <n-input v-model:value="form.Url" placeholder="https://www.vjshi.com/watch/8909441.html" />
              </n-form-item>
              <n-form-item path="VideoShort" label="预览视频">
                <n-input v-model:value="form.Payload.VideoShort" placeholder="https://lmp4.vjshi.com/2022-08-28/7cf436ce6cbb4b589dfae5a68318bd10.mp4" />
              </n-form-item>
              <n-form-item path="Pic" label="预览图片">
                <n-input v-model:value="form.Payload.Pic" placeholder="https://pic.vjshi.com/2022-08-28/7cf436ce6cbb4b589dfae5a68318bd10/online/797a9896d32e46958e8a3a0ac33a9052.jpg?x-oss-process=style/w1440_h2880" />
              </n-form-item>
              <n-form-item path="VideoPreview" label="视频小样">
                <n-input v-model:value="form.Payload.VideoPreview" placeholder="https://mp4.vjshi.com/2022-08-28/7cf436ce6cbb4b589dfae5a68318bd10.mp4" />
              </n-form-item>
              <n-form-item path="VideoFormal" label="视频地址/压缩包地址">
                <n-input v-model:value="form.Payload.VideoFormal" placeholder="https://down.vjshi.com/2022-08-28/8909441_VJshi_0480c6149f8fa56ab686590d33896492.mp4?auth_key=1701938201-0-0-c007876c45599831d2bb50d212021d6b&rename=%E5%85%89%E5%8E%82_8909441_%E6%BB%B4%E6%B0%B4%E7%A9%BF%E7%9F%B3%E5%B2%A9%E7%9F%B3%E6%BB%B4%E6%B0%B4.mp4" />
              </n-form-item>


              <n-form-item>
                <n-button v-if="form.Id === '' " type="primary" @click="createVideo"
                          :disabled="form.Payload.Number === '' ||
                          form.Payload.FormalLocal === '' ||
                          form.Payload.ShortLocal === '' ||
                          form.Payload.Title === ''"
                          :loading="loading"
                >创建</n-button>
                <n-button v-else type="primary" @click="createVideo">保存</n-button>
                <n-button @click="dialogVisible = false">取消</n-button>
              </n-form-item>
            </n-form>

          <template #footer>

          </template>
        </n-card>
      </n-modal>
    </el-main>
    </n-message-provider>
  </el-container>

</template>

<script>
import {defineEmits,defineComponent, h, ref,reactive, nextTick, onMounted, toRefs ,computed} from "vue";
import {ExportVideoList, VideoManage,VideoCreate,VideoDelete,VideoGetNumber } from "../../wailsjs/go/internal/App.js";
import {ElNotification} from "element-plus";

import { NTable, NDialog, NForm, NFormItem, NInput,useMessage,NTag, NButton } from 'naive-ui';

const emits = defineEmits(['emit-video-list']);

export default defineComponent({

  setup() {
    const message = useMessage()
    const rules = {
      Title: [
        {
          required: true,
          message: "请输入视频标题"
        }
      ],
      Number: [
        {
          required: true,
          message: "请输入视频编号"
        }
      ],
      FormalLocal: [
        {
          required: true,
          message: "请输入本地视频资源地址"
        }
      ],
      ShortLocal: [
        {
          required: true,
          message: "请输入本地预览视频"
        }
      ],
    };
    const keyword = ref()

    const options = [
      {
        label: "编辑",
        key: "edit",
      },
      {
        label: () => h("span", { style: { color: "red" } }, "删除"),
        key: "delete",
      }
    ];

    const paginationRef = computed(() => ({
      pageSize: 10,
    }))
    // const message = useMessage();
    const showDropdownRef = ref(false);
    const xRef = ref(0);
    const yRef = ref(0);

    const searchData = () =>  {
      VideoManage(keyword.value).then(res => {
        if (res.code !== 200) {
          ElNotification({
            title:res.msg,
            type:"error",
          })
        }
        keyword.value = res.data.Query

        if (res.data.Hits === 0){
          data.value = {}

        }else{
          data.value = res.data.Items
        }

        // console.log(data.value )
        // createData(res.data.Items)
      })
    }

    onMounted(() => {
      searchData()
    })

    let form = reactive({})
    let data = ref()

    const dialogVisible = ref(false)
    // const dialogVisibleAdd = ref(false)
    const modelRef = ref({
      Id: "",
      Url: "",
      Type: "",
      Payload:{
        Title:"",
        VideoShort: "",
        Pic: "",
        Url: "",
        VideoPreview: "",
        VideoFormal: "",
        Number: "",
        ShortLocal: "",
        PreviewLocal: "",
        FormalLocal: "",
        ImgLocal: "",
      }
    })

    const showDataModal = (row) =>{
      dialogVisible.value = true
      form.value = row
      modelRef.value = row
    }

    let videoManageList =  () =>{
       VideoManage(keyword.value).then(res => {
        if (res.code !== 200) {
          ElNotification({
            title:res.msg,
            type:"error",
          })
        }
        // data.value = res.data.Items
        // console.log(data.value )
        keyword.value = res.data.Query
         if (res.data.Hits === 0){
           data.value = {}
         }else{
           data.value = res.data.Items
         }
         // createData(res.data.Items)
      })
    }

    const deleteRow = (row) =>{
      console.log(row)

      VideoDelete(row).then(res => {
        dialogVisible.value = false
        if (res.code !== 200) {
          ElNotification({
            title:res.msg,
            type: "error",
          })
          return
        }
        ElNotification({
          title:res.msg,
          type: "success",
        })
        videoManageList()

      })
    }

    const addList = () =>{
      VideoCreate(modelRef.value).then(res => {
        dialogVisible.value = false
        if (res.code !== 200) {
          ElNotification({
            title:res.msg,
            type: "error",
          })
          return
        }
        ElNotification({
          title:res.msg,
          type: "success",
        })

      })
    }

    const GetNumber =() =>{
      VideoGetNumber(modelRef.value).then(res => {
        if (res.code !== 200) {
          ElNotification({
            title:res.msg,
            type: "error",
          })
          return
        }
        console.log(res.data)
        modelRef.value.Payload.Number = res.data.Payload.Number

      })

    }

    return {
      emits,
      showDataModal,
      deleteRow,
      videoManageList,
      addList,
      dialogVisible,
      rules,
      form :modelRef,
      data,
      keyword,
      paginationRef,

      searchData,
      options,
      showDropdown: showDropdownRef,
      x: xRef,
      y: yRef,
      selectedRow: {},
      rowKey (rowData) {
        return rowData
      },

      createVideo() {
        addList()
        setTimeout(() => {
          // 在这里编写需要延迟执行的代码
          videoManageList()
        }, 100); // 1000 是延迟1秒执行

      },

      dialogVisibleAdd(){
        modelRef.value = ({
          Id: "",
          Url: "",
          Type: "",
          Payload:{
            Title:"",
            VideoShort: "",
            Pic: "",
            Url: "",
            VideoPreview: "",
            VideoFormal: "",
            Number: "",
            ShortLocal: "",
            PreviewLocal: "",
            FormalLocal: "",
            ImgLocal: "",
          }
        })
        dialogVisible.value = true
      },

      VideoListExport(){
        ExportVideoList(keyword.value).then(res => {
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
        videoManageList()
      },

      handleSelect() {
        showDropdownRef.value = false;
      },

      onClickoutside() {
        showDropdownRef.value = false;
      },

      rowProps: (row) => {
        return {
          onContextmenu: (e) => {
            // message.info(JSON.stringify(row, null, 2));
            // alert(JSON.stringify(row, null, 2))

            e.preventDefault();
            showDropdownRef.value = false;
            nextTick().then(() => {
              showDropdownRef.value = true;
              xRef.value = e.clientX;
              yRef.value = e.clientY;
              // dialogVisible = true;
              dialogVisible.value = false

            });
          }
        };
      },

      cols : [
        {
          title: "No.",
          key: "no"
        },
        {
          title: "id",
          key: "Id"
        },
        {
          title: "标题",
          key: "Title",
          render (row) {
            return h(
                NButton,
                {
                  size: "small",
                  onClick: () => showDataModal(row)
                },
                { default: () => row.Payload.Title}
            );
          }
        },
        {
          title: "操作",
          key: "Title",
          render (row) {
            return h(
                NButton,
                {
                  size: "small",
                  onClick: () => deleteRow(row)
                },
                { default: () => "删除"}
            );
          }
        },
        {
          title: "网站链接",
          key: "Url"
        },
        // {
        //   title: '本地链接',
        //   key: 'ShortLocal',
        // }
      ],
      autoGetVideoNumber(){
        message.info("待开发")
        GetNumber()
      }
    };
  }
});
</script>