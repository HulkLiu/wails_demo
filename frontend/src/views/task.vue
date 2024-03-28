<template>
  <el-container class="page-cantor" v-loading="loading" :style="areaStyle">
    <el-header class="header-area" style="height: 89px;">
      <n-button @click="dialogVisibleAdd"> 新增任务 </n-button>
    </el-header>
    <n-message-provider>

      <el-main class="main-area" style="--wails-draggable:no-drag">

        <n-data-table
            :columns="columns"
            :data="listData"
            :pagination="pagination"
            :bordered="false"
        />

        <n-modal v-model:show="dialogVisible" style="width: 1500px">
          <n-card
              style="width: 67%"
              title="任务新增/修改"
              :bordered="false"
              size="huge"
              role="dialog"
              aria-modal="true"

              :style="{maxWidth: '840px'}"
              :label-width="90"
          >

            <template #header-extra></template>

            <n-form  :model="form"  label-placement="left">
              <n-form-item v-show="false" path="id" label="id">
                <n-input :disabled="true" v-model:value="form.Id"  />
              </n-form-item>

              <n-form-item path="name" label="任务名称">
                <n-input v-model:value="form.Name" placeholder="任务名称" />
              </n-form-item>
              <n-form-item path="Description" label="任务描述">
                <n-input v-model:value="form.Description" placeholder="任务描述" />
              </n-form-item>
              <n-form-item label="完成情况" path="completed">
                <n-select
                    v-model:value="form.Completed"
                    placeholder="Select"
                    :options="generalOptions"
                />
              </n-form-item>
              <n-form-item path="createdAt" label="创建时间">
                <n-input :disabled="true" v-model:value="form.CreatedAt" placeholder="创建时间" />
              </n-form-item>
              <n-form-item path="updatedAt" label="更新时间">
                <n-input :disabled="true" v-model:value="form.UpdatedAt" placeholder="更新时间" />
              </n-form-item>

              <n-form-item>
                <n-button v-if="form.Id === '' " type="primary" @click="onTaskAdd"
                          :disabled="form.Name === '' ||
                          form.Description === '' ||
                          form.Completed === ''"
                          :loading="loading"
                >创建</n-button>
                <n-button v-else type="primary" @click="onTaskAdd">保存</n-button>
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
import {h, defineComponent, onMounted, ref, reactive} from "vue";
import { NButton, useMessage } from "naive-ui";
import {AddTask, GetTaskList,UpdateTask} from "../../wailsjs/go/internal/App.js";
import moment from 'moment';
import {ElNotification} from "element-plus";



export default defineComponent({

  setup() {
    const listData = ref()
    // const listData = reactive({})

    const message = useMessage()
    const form = ref()

    const taskList = () =>{
      GetTaskList().then(res => {
        console.log(res)

        if (res.code === 0){
          if (res.data === null){
            listData.value = []
          }else{
            listData.value = res.data;

            // listData.value.CreatedAt = moment(listData.value.CreatedAt).format("YYYY-MM-DD HH:mm:ss")
            // listData.value.UpdatedAt = moment(listData.value.UpdatedAt).format("YYYY-MM-DD HH:mm:ss")

          }
          console.log(listData.value)
        }
      })
    }

    const taskAdd = () =>{
      if (modelRef.value.CreatedAt==="" ){
        modelRef.value.CreatedAt = new Date().toISOString();
      }
      if ( modelRef.value.UpdatedAt===""){
        modelRef.value.UpdatedAt = new Date().toISOString();
      }
      console.log(modelRef.value.Id)
      if (modelRef.value.Id === null) {
        modelRef.value.Id = 0
        AddTask(JSON.stringify(modelRef.value) ).then(res => {
          console.log(res)
          if (res.code === 0) {
            message.success("添加成功")
            dialogVisible.value = false
            taskList()

          }else{
            message.warning("添加失败")
            message.error(res.msg)
          }

        })
      }else{
        editTask()
      }


    }

    const showEditTask = (row) =>{
      dialogVisible.value = true
      modelRef.value = row
      console.log(modelRef.value )
    }
    const editTask = (row) =>{
      UpdateTask(JSON.stringify(modelRef.value) ).then(res => {
        console.log(res)
        if (res.code === 0) {
          message.success("编辑成功")
          dialogVisible.value = false

        }else{
          message.warning("编辑失败")
          message.error(res.msg)
        }

        taskList()
      })


    }

    onMounted(() => {
      taskList()
    })
    const dialogVisible = ref(false)
    const modelRef = ref({
      Id:null ,
      Name:"Learn Go ",
      Description:"Study the basics of Go" ,
      Completed:0 ,
      CreatedAt:"2023-12-28 17:36:47.7258618 +0800 CST m=+16.554803701 ",
      UpdatedAt:"2023-12-28 17:36:47.7258618 +0800 CST m=+16.554803701",
    })
    return {
      dialogVisible,
      listData,
      form:modelRef,
      pagination: { pageSize: 10 },
      taskList,
      onTaskAdd(){
        taskAdd()
      } ,

      dialogVisibleAdd(){
        modelRef.value = ({
          Id:null ,
          Name:"" ,
          Description:"" ,
          Completed:0,
          CreatedAt:"",
          UpdatedAt:"",
        })
        dialogVisible.value = true
      },
      generalOptions: [0,1].map(
          (v) => ({
            label: v,
            value: v
          })
      ),
      columns : [
        {
          type: "selection",
        },
        {
          title: "任务编号",
          key: "ID"
        },
        {
          title: "完成度",
          key: "Completed",
          defaultFilterOptionValues: [],
          filterOptions: [
            {
              label: "完成",
              value:1,

            },
            {
              label: "未完成",
              value: 0
            }
          ],
          filter(value, row) {
            return !!~row.Completed.indexOf(value);
          },
          render (row) {
            return h(
                NButton,
                {
                  size: "small",
                  onClick: () => showEditTask(row)
                },
                { default: () => row.Completed}
            );
          }
        },
        {
          title: "任务名称",
          key: "Name"
        },
        {
          title: "任务描述",
          key: "Description"
        },

        {
          title: "创建时间",
          key: "CreatedAt",

        },
        {
          title: "修改时间",
          key: "UpdatedAt"
        },
        // {
        //   title: "操作",
        //   key: "Title",
        //   render (row) {
        //     return h(
        //         NButton,
        //         {
        //           size: "small",
        //           onClick: () => alert(row)
        //         },
        //         { default: () => "删除"}
        //     );
        //   }
        // },
      ],


    };
  }
});
</script>