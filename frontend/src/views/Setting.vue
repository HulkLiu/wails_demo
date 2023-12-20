<template>
  <div>
    <n-card :bordered="false" style="height: 100%;">
      <n-tabs type="line" animated default-value="system">
        <n-tab-pane name="system" tab="软件设置">
          <n-space vertical>
            <n-card title="软件优化" class="card-radius-10" embedded>
              <n-collapse :default-expanded-names="['2']">
                <n-collapse-item title="前端配置" name="1">

                </n-collapse-item>
                <n-collapse-item title="后端配置" name="2">
                  <n-space vertical>
                    <n-switch v-model:value="updateDisabledBak" />
                    <n-form
                        ref="formRef"
                        :model="modelBak"
                        label-placement="left"
                        :label-width="160"
                        :disabled="updateDisabledBak"
                        :style="{maxWidth: '640px'}"
                    >
                      <n-form-item label="ES 仓库 Index" path="inputValue">
                        <n-input v-model:value="modelBak.Bak.EsIndex" placeholder="Es 仓库索引" />
                        <n-button attr-type="button" @click="updateModelBak" type="success">
                          更新
                        </n-button>
                      </n-form-item>

                      <n-form-item label="本地文件配置路径" path="inputValue">
                        <n-input v-model:value="modelBak.Bak.LocalPath" placeholder="本地文件配置路径" />
                        <n-button attr-type="button" @click="updateModelBak" type="success">
                          更新
                        </n-button>
                      </n-form-item>

                      <n-form-item label="Excel保存目录" path="inputValue">
                        <n-input v-model:value="modelBak.Bak.ExcelDir" placeholder="Excel保存目录" />
                        <n-button attr-type="button" @click="updateModelBak" type="success">
                          更新
                        </n-button>
                      </n-form-item>

                      <n-form-item label="ES容器ID" path="inputValue">
                        <n-input v-model:value="modelBak.Bak.ContainerID" placeholder="ES容器ID" />
                        <n-button attr-type="button" @click="updateModelBak" type="success">
                          更新
                        </n-button>
                      </n-form-item>

                    </n-form>
                  </n-space>
                </n-collapse-item>
                <n-collapse-item title="其他配置" name="3">
                  <n-space vertical>
                    <n-switch v-model:value="updateDisabled" />
                    <n-form
                        ref="formRef"
                        :model="model"
                        label-placement="left"
                        :label-width="160"
                        :disabled="updateDisabled"
                        :style="{maxWidth: '640px'}"
                    >
                      <n-form-item label="Input" path="inputValue">
                        <n-input v-model:value="model.inputValue" placeholder="Input" />
                      </n-form-item>
                      <n-form-item label="Select" path="selectValue">
                        <n-select
                            v-model:value="model.selectValue"
                            placeholder="Select"
                            :options="generalOptions"
                        />
                      </n-form-item>
                      <n-form-item label="Upload" path="uploadValue">
                        <n-upload>
                          <n-button>Upload file</n-button>
                        </n-upload>
                      </n-form-item>
                      <n-form-item label="Auto Complete" path="autoCompleteValue">
                        <n-auto-complete
                            v-model:value="model.autoCompleteValue"
                            :options="autoCompleteOptions"
                            placeholder="Email"
                        />
                      </n-form-item>
                      <n-form-item label="Dynamic Tags" path="dynamicTagsValue">
                        <n-dynamic-tags v-model:value="model.dynamicTagsValue" />
                      </n-form-item>
                      <n-form-item label="Tree Select" path="treeSelectValue">
                        <n-tree-select
                            :options="treeSelectOptions"
                            default-value="Drive My Car"
                        />
                      </n-form-item>
                      <n-form-item label="Cascader" path="cascaderValue">
                        <n-cascader
                            v-model:value="model.cascaderValue"
                            placeholder="Cascader"
                            :options="options"
                            check-strategy="child"
                            size="medium"
                        />
                      </n-form-item>
                      <n-form-item label="Datetime" path="datetimeValue">
                        <n-date-picker v-model:value="model.datetimeValue" type="datetime" />
                      </n-form-item>
                      <n-form-item label="Switch" path="switchValue">
                        <n-switch v-model:value="model.switchValue" />
                      </n-form-item>
                      <n-form-item label="Checkbox" path="checkboxValue">
                        <n-checkbox v-model:checked="model.checkboxValue">
                          Checkbox
                        </n-checkbox>
                      </n-form-item>
                      <n-form-item label="Checkbox Group" path="checkboxGroupValue">
                        <n-checkbox-group v-model:value="model.checkboxGroupValue">
                          <n-space>
                            <n-checkbox value="Option 1">
                              Option 1
                            </n-checkbox>
                            <n-checkbox value="Option 2">
                              Option 2
                            </n-checkbox>
                            <n-checkbox value="Option 3">
                              Option 3
                            </n-checkbox>
                          </n-space>
                        </n-checkbox-group>
                      </n-form-item>
                      <n-form-item label="Radio" path="radioValue">
                        <n-radio
                            :checked="model.radioValue === 'Definitely Maybe'"
                            value="Definitely Maybe"
                            name="basic-demo"
                        >
                          Definitely Maybe
                        </n-radio>
                      </n-form-item>
                      <n-form-item label="Radio Group" path="radioGroupValue">
                        <n-radio-group v-model:value="model.radioGroupValue" name="radiogroup2">
                          <n-radio value="Radio 1">
                            Radio 1
                          </n-radio>
                          <n-radio value="Radio 2">
                            Radio 2
                          </n-radio>
                          <n-radio value="Radio 3">
                            Radio 3
                          </n-radio>
                        </n-radio-group>
                      </n-form-item>
                      <n-form-item label="Radio Button Group" path="radioGroupValue">
                        <n-radio-group v-model:value="model.radioGroupValue" name="radiogroup2">
                          <n-radio-button value="Radio 1">
                            Radio 1
                          </n-radio-button>
                          <n-radio-button value="Radio 2">
                            Radio 2
                          </n-radio-button>
                          <n-radio-button value="Radio 3">
                            Radio 3
                          </n-radio-button>
                        </n-radio-group>
                      </n-form-item>
                      <n-form-item label="Input Number" path="inputNumberValue">
                        <n-input-number v-model:value="model.inputNumberValue" />
                      </n-form-item>
                      <n-form-item label="Time Picker" path="timePickerValue">
                        <n-time-picker v-model:value="model.timePickerValue" />
                      </n-form-item>
                      <n-form-item label="Color Picker" path="colorValue">
                        <n-color-picker v-model:value="model.colorValue" />
                      </n-form-item>
                      <n-form-item label="Slider" path="sliderValue">
                        <n-slider v-model:value="model.sliderValue" :step="5" />
                      </n-form-item>
                      <n-form-item label="Transfer" path="transferValue">
                        <n-transfer
                            v-model:value="model.transferValue"
                            :options="generalOptions"
                        />
                      </n-form-item>
                    </n-form>
                  </n-space>

                </n-collapse-item>
              </n-collapse>
            </n-card>
          </n-space>
        </n-tab-pane>
        <n-tab-pane name="imgbed" tab="图床配置">
          <n-space vertical>
            <n-card title="本地存储配置" embedded class="card-radius-10">
            </n-card>
            <n-card title="阿里云OSS配置" embedded class="card-radius-10">
              <n-alert title="注意事项" type="info" v-show="data.imgbed.configType == 'alioss'" :bordered="false">
                Region、Bucket、Object内容请勿添加任何符号
              </n-alert>
            </n-card>
          </n-space>
        </n-tab-pane>
        <n-tab-pane name="ocr" tab="图文配置">
          <n-space vertical>
            <n-card title="百度OCR" embedded class="card-radius-10">
            </n-card>
            <n-card title="百度翻译(翻译开放平台)" embedded class="card-radius-10">
              <n-alert title="注意事项" type="info" :bordered="false">
                百度翻译免费额度下并发量为1，请勿在短时间内快速、重复地进行翻译。<br>
                appid 和 secret 务必准确，其他三项根据实际需求配置。
              </n-alert>
            </n-card>
            <n-card title="腾讯云翻译" embedded class="card-radius-10">
              <n-alert title="注意事项" type="info" :bordered="false">
                腾讯云翻译免费额度下默认接口请求频率限制：5次/秒，文本翻译的每月免费额度为5百万字符，免费额度每月1日0点重置。<br>
              </n-alert>
            </n-card>
            <n-card title="阿里云机器翻译" embedded class="card-radius-10">
              <n-alert title="注意事项" type="info" :bordered="false">
                阿里云机器翻译免费额度下QPS为50，通用版机器翻译每月100万字符免费额度，单次请求字符不要太大(控制在1000字符以内)，否则请使用文本翻译。<br>
              </n-alert>
              待对接
            </n-card>
          </n-space>
        </n-tab-pane>
        <n-tab-pane name="other" tab="其他配置">
          待开发
        </n-tab-pane>
        <n-tab-pane name="instructions" tab="使用说明">
          待开发
        </n-tab-pane>
        <n-tab-pane name="feedback" tab="问题反馈">
          待开发
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>


</template>

<script>
import {defineComponent, computed, ref, onMounted} from "vue";
import { useMessage } from 'naive-ui'
import {GetSettingList,UpdateSettingList} from "../../wailsjs/go/internal/App.js";
import {ElNotification} from "element-plus";

function genOptions(depth = 2, iterator = 1, prefix = "") {
  const length = 12;
  const options = [];
  for (let i = 1; i <= length; ++i) {
    if (iterator === 1) {
      options.push({
        value: `${i}`,
        label: `${i}`,
        disabled: i % 5 === 0,
        children: genOptions(depth, iterator + 1, "" + i)
      });
    } else if (iterator === depth) {
      options.push({
        value: `${prefix}-${i}`,
        label: `${prefix}-${i}`,
        disabled: i % 5 === 0
      });
    } else {
      options.push({
        value: `${prefix}-${i}`,
        label: `${prefix}-${i}`,
        disabled: i % 5 === 0,
        children: genOptions(depth, iterator + 1, `${prefix}-${i}`)
      });
    }
  }
  return options;
}
export default defineComponent({
  setup() {
    const data = ref({
      // 默认值必须有，否则会重复触发更新方法
      imgbed: {
        configType: null
      },
      trans: {
        transType: null
      }
    })
    const formRef = ref(null);
    const formRefBak = ref(null);
    const model = ref({
      inputValue: null,
      selectValue: null,
      autoCompleteValue: "",
      dynamicTagsValue: ["teacher", "frontend"],
      cascaderValue: null,
      datetimeValue: null,
      switchValue: false,
      checkboxValue: false,
      checkboxGroupValue: null,
      radioValue: "Definitely Maybe",
      radioGroupValue: null,
      radioButtonGroupValue: null,
      inputNumberValue: null,
      timePickerValue: null,
      colorValue: null,
      sliderValue: 0,
      transferValue: null
    });
    const modelBak = ref({
      Bak:{
        EsIndex: "",
        LocalPath: "",
        ContainerID: "",
        ExcelDir: "",
      }
    });
    const bakData = () => {
      GetSettingList().then(res => {
        if (res.code !== 200) {
          ElNotification({
            title: res.msg,
            type: "error",
          })
          return
        }
        // console.log(res.data)
        modelBak.value = res.data
        // console.log(modelBak.value)
      })
    }
    const upBak = () => {
      UpdateSettingList(modelBak.value).then(res => {
        if (res.code !== 200) {
          ElNotification({
            title: res.msg,
            type: "error",
          })
          return
        }
        message.success("更新成功")
      })
    }

    onMounted(() => {
      bakData()
    })
    const message = useMessage()
    return {
      data,
      upBak,
      bakData,
      updateDisabled: ref(false),
      updateDisabledBak: ref(false),
      formRef,
      formRefBak,
      model,
      modelBak,
      updateModelBak(){
        upBak()

        setTimeout(() => {
          // 在这里编写需要延迟执行的代码
          bakData()
        }, 1000); // 延迟1秒执行
        // alert(modelBak.value)
      },
      generalOptions: ["groode", "veli good", "emazing", "lidiculous"].map(
          (v) => ({
            label: v,
            value: v
          })
      ),
      options: genOptions(),
      treeSelectOptions: [
        {
          label: "Rubber Soul",
          key: "Rubber Soul",
          children: [
            {
              label: "Drive My Car",
              key: "Drive My Car"
            },
            {
              label: "Michelle",
              key: "Michelle"
            }
          ]
        }
      ],
      autoCompleteOptions: computed(() => {
        return ["@gmail.com", "@163.com", "@qq.com"].map((suffix) => {
          const prefix = model.value.autoCompleteValue.split("@")[0];
          return {
            label: prefix + suffix,
            value: prefix + suffix
          };
        });
      })
    };
  }
});


</script>

<style scoped>
.card-radius-10 {
  border-radius: 10px;
}
</style>
