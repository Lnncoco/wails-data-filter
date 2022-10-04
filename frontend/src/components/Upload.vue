<template>
  <div class="upload-container" data-wails-no-drag>
    <div class="hero-overlay bg-opacity-0"></div>
    <div class="hero-content text-center text-neutral-content">
      <div class="max-w-md">
        <div class="file-info-area mt-4 h-56">
          <div v-show="state === STATE_INIT" class="pt-4">
            <h1 class="mb-5 text-5xl font-bold text-slate-900 select-none">
              数据过滤工具
            </h1>
            <p v-show="!filePath" class="mb-5 text-gray-400 select-none">
              请选择目标文件，程序会根据配置处理相应字符
            </p>
          </div>
          <div class="progress-area" v-show="state !== STATE_INIT">
            <div class="stats shadow w-full">
              <div class="stat text-left">
                <div class="stat-title truncate">{{ filename }}</div>
                <div
                  v-if="lineNum && state === STATE_WAIT"
                  class="stat-value mt-2 text-center"
                >
                  <span class="stat-desc font-normal">已处理</span>
                  {{ lineNum }}
                  <span class="stat-desc font-normal">行</span>
                </div>
                <div
                  v-else-if="state === STATE_DONE"
                  class="stat-value mt-2 text-center"
                >
                  <span class="stat-desc text-2xl">处理完成</span>
                </div>
                <div v-else class="stat-value mt-2 text-center">
                  <span class="stat-desc text-2xl">等待处理</span>
                </div>
                <div class="stat-desc mt-2">{{ filePath }}</div>
                <div
                  v-show="state === STATE_DONE"
                  class="stat-desc truncate mt-2"
                >
                  新文件：{{ newFilePath }}
                </div>
                <div
                  v-show="state === STATE_DONE"
                  class="stat-desc truncate mt-2"
                >
                  耗时：{{ elapsedTime }}
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          class="select-area relative transition duration-500 select-none cursor-pointer mx-auto mt-1"
          :class="{
            'state-init': state === STATE_INIT,
            'state-ready': state === STATE_READY,
            'state-wait': state === STATE_WAIT,
            'state-done': state === STATE_DONE,
          }"
          @click="handleClickButton"
        >
          <div class="add-files-area flex transition duration-300 z-10">
            <div class="add-files button w-full">
              <span class="file-ico ico-add"
                ><svg-icon name="file-medical"
              /></span>
              <span class="file-ico ico-exist"
                ><svg-icon name="file-alt"
              /></span>
              <span class="file-ico ico-wait"><svg-icon name="loading" /></span>
              <span class="file-ico ico-success"
                ><svg-icon name="file-check-alt"
              /></span>
              <span>{{ buttonText }}</span>
            </div>
          </div>
          <div class="reupload-area absolute" @click.stop="handleSelectFile">
            <div class="button text-center absolute inset-0">重新选择</div>
          </div>
        </div>
      </div>
    </div>
    <input type="file" ref="fileEleRef" id="fileEle" class="hidden" />
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from "vue";
import {
  SelectFile,
  FilterFile,
  SplitFilePath,
  EventsOn,
  LogInfo,
} from "../../wailsjs";

// 状态 1 等待选择  2 已选择文件 3 处理中 4 已处理
const STATE_INIT = 1;
const STATE_READY = 2;
const STATE_WAIT = 3;
const STATE_DONE = 4;
const state = ref(STATE_INIT);

const buttonText = ref("选择文件");
const filePath = ref("");
const filename = ref("");
const newFilePath = ref("");
const lineNum = ref(0);
const elapsedTime = ref(0);

const handleSelectFile = async () => {
  const path = await SelectFile();
  LogInfo(`选择的文件路径为：${path}`);
  if (path) {
    state.value = STATE_READY;
    filePath.value = path;
    SplitFilePath(path).then((res) => {
      filename.value = `${res.filename}.${res.ext}`;
      newFilePath.value = `${res.filename}${res.addSuffixStr}.${res.ext}`;
    });
  }
};

const handleClickButton = () => {
  switch (state.value) {
    case STATE_INIT:
      handleSelectFile();
      break;
    case STATE_READY:
      filterFile(filePath.value);
      break;
    case STATE_DONE:
      state.value = STATE_INIT;
      filePath.value = "";
      newFilePath.value = "";
      lineNum.value = 0;
      elapsedTime.value = 0;
      handleSelectFile();
      break;
  }
};

const filterFile = async (path) => {
  state.value = STATE_WAIT;
  const time = await FilterFile(path);
  state.value = STATE_DONE;
  LogInfo(`处理结果：${time}`);
  if (time) {
    elapsedTime.value = time;
  }
};

const onFilterChange = () => {
  EventsOn("filter-change", (num) => {
    lineNum.value = num;
  });
};

watch(state, (val) => {
  if (val === STATE_INIT) buttonText.value = "选择文件";
  else if (val === STATE_READY) buttonText.value = "处理文件";
  else if (val === STATE_WAIT) buttonText.value = "处理中...";
  else if (val === STATE_DONE) {
    buttonText.value = "处理完成";
  }
});

onMounted(() => {
  onFilterChange();
});
</script>

<style lang="less" scoped>
@import "../animation.less";
.upload-container {
  .text-focus-in();
  .file-info-area {
    .progress-area {
      .stats {
        background-color: transparent;
        box-shadow: 0 0 16px 1px rgb(0 0 0 / 20%);
      }
    }
  }
  .select-area {
    width: 240px;
    height: 60px;
    border-radius: 2.5rem;
    color: #252525;
    box-shadow: 0 0 16px 1px rgb(0 0 0 / 25%);

    &.state-init {
      .ico-add {
        display: block;
      }
      .ico-exist,
      .ico-wait,
      .ico-success {
        display: none;
      }
    }
    &.state-ready {
      width: 326px;
      .ico-exist,
      .reupload-area {
        display: block;
      }
      .ico-add,
      .ico-wait,
      .ico-success {
        display: none;
      }
    }
    &.state-wait {
      .ico-wait {
        display: block;
        .ele-rotate360();
      }
      .ico-add,
      .ico-exist,
      .ico-success {
        display: none;
      }
    }
    &.state-done {
      .ico-success {
        display: block;
      }
      .ico-add,
      .ico-wait,
      .ico-exist {
        display: none;
      }
    }

    .add-files-area {
      padding-left: 35px;
      border-radius: 2.5rem;
      &:hover {
        background-color: #fdda65;
      }
    }
    .add-files {
      font-size: 2rem;
      font-weight: bolder;
      display: flex;
      flex-flow: row;
      align-items: center;
      span {
        font-size: 25px;
        line-height: 60px;
      }
    }

    .reupload-area {
      display: none;
      font-size: 0.875rem;
      background: #f1f1f1;
      height: 3rem;
      border-radius: 1.875rem;
      padding: 0.41667rem;
      transition: 0.25s;
      width: 6.75rem;
      top: 0.4rem;
      right: 0.83333rem;
      &:hover {
        background-color: #fdda65;
      }
      .button {
        width: 6.75rem;
        font-size: 1rem;
        font-weight: 700;
        height: 3rem;
        border-radius: 1.875rem;
        line-height: 3rem;
      }
    }
    .add-files .file-ico {
      font-size: 2.2rem;
      line-height: 60px;
      margin-right: 17px;
      transform: translateY(-1px);
    }
  }
}
</style>
