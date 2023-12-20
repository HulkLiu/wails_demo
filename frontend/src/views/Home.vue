<template>
  <n-space vertical size="large">
    <n-layout>
      <n-layout-header style="margin-top: 12px">基本统计</n-layout-header>
      <n-layout-content  content-style="padding: 24px;">
        <n-row>
          <n-col :span="12">
            <n-statistic  label="视频资源数量" :value="refData.value">
              <template #prefix>
                <n-icon>
                  <md-save />
                </n-icon>
              </template>

              <template #suffix>
                <n-popover trigger="hover">
                  <template #trigger>
                    <n-span style="color: #3a51bb" >{{refData.videoCount}} 个</n-span>
                  </template>
                  <div v-for="(item, index) in refData.videoType"  >
                    <span style="font-size: 20px;">{{ index }} : {{ item }}</span>
                  </div>
                </n-popover>
              </template>
            </n-statistic>
          </n-col>
          <n-col :span="12">
            <n-statistic label="活跃用户">

              <n-popover trigger="hover">
                <template #trigger>
                  <n-span style="color: #3a51bb">{{refData.userInfo.length}} 位</n-span>
                </template>
                <div v-for="(item, index) in refData.userInfo"  >
                  <span style="font-size: 20px">{{ index +1  }} : {{ item.user_repo }} - {{ item.file_owner }}</span>
                </div>
              </n-popover>
            </n-statistic>
          </n-col>
        </n-row>

      </n-layout-content>
      <n-layout-footer >其他统计</n-layout-footer>
      <div id="chartContainer" style="width: 800px;height:400px;margin-top: 27px"></div>
    </n-layout>
  </n-space>
</template>
<script >
import {defineComponent, onMounted, reactive, ref} from "vue";
import { MdSave } from "@vicons/ionicons4";
import {GetHomeInfo} from "../../wailsjs/go/internal/App.js";
import {ElNotification} from "element-plus";
import * as echarts from 'echarts'



export default defineComponent({
  components: {
    MdSave
  },

  setup(){
    let list = ref()
    const refData = ref({
      userInfo:"",
      videoCount:"",
      videoType:"",
      dirInfo:{
        filed:[],
        count:[],
        size:[],
      },
    })

    const HomeInfo = () =>{
      GetHomeInfo().then(res => {
        if (res.code !== 200) {
          ElNotification({
            title:res.msg,
            type:"error",
          })
        }
        console.log(res.data)

        refData.value.userInfo = res.data.User
        refData.value.videoCount = res.data.Video.Video
        refData.value.videoType = res.data.Video.Type

        refData.value.dirInfo.filed = res.data.Dir.Filed
        refData.value.dirInfo.count = res.data.Dir.Count
        refData.value.dirInfo.size = res.data.Dir.Size
        console.log(refData.value.dirInfo.filed )
        initChart()

        // refData.value.dirInfo.filed = JSON.stringify(res.data.Dir.Filed).replace(/\\"/g, '').replace(/^"|"$/g, '')
        // refData.value.dirInfo.count = JSON.stringify(res.data.Dir.Count).replace(/^"|"$/g, '')
        // refData.value.dirInfo.size = JSON.stringify(res.data.Dir.Size).replace(/^"|"$/g, '')

        // console.log(JSON.stringify(refData.value.dirInfo.count).replace(/^"|"$/g, ''))
        // console.log(JSON.stringify(refData.value.dirInfo.size ).replace(/^"|"$/g, ''))
        // console.log(JSON.stringify(refData.value.dirInfo.filed ).replace(/\\"/g, '').replace(/^"|"$/g, ''))
        // console.log([2.0, 2.2, 3.3, 4.5, 6.3, 10.2, 20.3, 23.4, 23.0, 16.5, 12.0, 6.2])
        // createData(res.data.Items)
      })

    }

    const initChart = ()  =>{
      const chartContainer = document.getElementById('chartContainer')
      const chart = echarts.init(chartContainer)
      // 在这里配置你的图表选项和数据
      // 例如：chart.setOption({...})
      chart.setOption({
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            crossStyle: {
              color: '#999'
            }
          }
        },
        toolbox: {
          feature: {
            dataView: { show: true, readOnly: false },
            magicType: { show: true, type: ['line', 'bar'] },
            restore: { show: true },
            saveAsImage: { show: true }
          }
        },
        legend: {
          data: [refData.value.dirInfo.filed,refData.value.dirInfo.count,refData.value.dirInfo.size]
        },
        xAxis: [
          {
            type: 'category',
            data: refData.value.dirInfo.filed,
            axisPointer: {
              type: 'shadow'
            }
          }
        ],
        yAxis: [
          {
            type: 'value',
            name: '文件数量',
            min: 0,
            max: getMax(refData.value.dirInfo.count) + getAvg(refData.value.dirInfo.count),
            interval: Math.ceil( getMax(refData.value.dirInfo.count)/3 ),
            axisLabel: {
              formatter: '{value} 个'
            }
          },
          {
            type: 'value',
            name: '文件大小',
            min: 0,
            max: getMax(refData.value.dirInfo.size) ,
            interval: getMax(refData.value.dirInfo.size)/ 5 ,
            axisLabel: {
              formatter: '{value} MB'
            }
          }
        ],
        series: [
          {
            name: '文件数',
            type: 'bar',
            tooltip: {
              valueFormatter: function (value) {
                return value + ' 个';
              }
            },
            data: refData.value.dirInfo.count,

          },
          {
            name: 'Size ',
            type: 'line',
            yAxisIndex: 1,
            tooltip: {
              valueFormatter: function (value) {
                return value + ' MB';
              }
            },
            data: refData.value.dirInfo.size,
          }
        ]
      });
    }

    onMounted(() => {
      HomeInfo()
    })
    return{
      list,
      refData,
      initChart,

    }

  }


});

// 计算最大值
function getMax(numbers) {
  return Math.ceil(Math.max(...numbers));
}

// 计算最小值
function getMin(numbers) {
  return Math.ceil(Math.min(...numbers));
}

// 计算平均值
function getAvg(numbers) {
  const sum = numbers.reduce((acc, curr) => acc + curr, 0);
  return Math.ceil(sum / numbers.length);
}

// 计算中位数
function getMedian(numbers) {
  numbers.sort((a, b) => a - b);
  const n = numbers.length;
  if (n % 2 === 0) {
    return Math.ceil((numbers[n / 2 - 1] + numbers[n / 2]) / 2);
  } else {
    return Math.ceil(numbers[Math.floor(n / 2)]);
  }
}

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
</style>
