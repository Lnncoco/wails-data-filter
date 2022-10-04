<template>
  <div class="menu-area absolute right-5 top-4" data-wails-no-drag>
    <ul class="menu bg-base-100 p-2 rounded-box text-xl">
      <li>
        <a @click="handleOptions" class="p-2.5">
          <svg-icon name="setting-config" />
        </a>
      </li>
      <li>
        <a @click="handleQuit" class="p-2.5">
          <svg-icon name="power" />
        </a>
      </li>
    </ul>
  </div>

  <div
    class="options-panel absolute bg-white px-5 py-7 w-7/12 h-5/6 rounded-2xl"
    v-show="showOptionsPanel"
    data-wails-no-drag
  >
    <div class="form-control">
      <label
        class="label cursor-pointer px-8 py-2 rounded-xl hover:bg-slate-50"
      >
        <span class="label-text text-slate-900 text-lg">清除零宽字符</span>
        <input
          type="checkbox"
          :checked="config.eraseZeroWidthCharacter && 'checked'"
          class="checkbox checkbox-primary"
        />
      </label>
      <!-- <label class="label cursor-pointer px-8 py-2 rounded-xl hover:bg-slate-50">
        <span class="label-text text-slate-900 text-lg">转义HTML实体符</span>
        <input
          type="checkbox"
          :checked="config.decodeHtmlEntity && 'checked'"
          class="checkbox checkbox-primary"
        />
      </label> -->
      <div class="divider">自定义替换规则</div>
      <div
        class="w-full px-8 py-2 pb-4 rounded-xl hover:bg-slate-50"
        v-for="(item, i) in config.regex"
      >
        <label class="label">
          <span class="label-text">正则表达式文本</span>
        </label>
        <input
          type="text"
          placeholder="正则表达式"
          v-model="config.regex[i].regexp"
          class="input input-bordered w-full"
        />
        <label class="label">
          <span class="label-text">替换的字符串</span>
        </label>
        <input
          type="text"
          placeholder="替换的字符串"
          v-model="config.regex[i].value"
          class="input input-bordered w-full"
        />
      </div>
      <button
        class="btn btn-wide btn-outline mx-auto mt-8 confirm"
        @click="confirm"
      >
        确定
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from "vue";
import { Quit, GetOptions, SetOptions, LogInfo } from "../../wailsjs";

const config = reactive({
  eraseZeroWidthCharacter: true,
  decodeHtmlEntity: false,
  regex: [],
});
const showOptionsPanel = ref(false);

const handleQuit = () => Quit();

const initRegex = (info) => {
  if (info && info.length) return info;
  return [{ regexp: "", value: "" }];
};

const handleOptions = async () => {
  const res = await GetOptions();
  LogInfo(`获取到的配置：${JSON.stringify(res)}`);
  console.log("获取到的配置：", res);
  config.eraseZeroWidthCharacter = res.eraseZeroWidthCharacter;
  config.decodeHtmlEntity = res.decodeHtmlEntity;
  config.regex = initRegex(res.regex);
  showOptionsPanel.value = true;
};

const confirm = () => {
  showOptionsPanel.value = false;
  const data = JSON.stringify({
    eraseZeroWidthCharacter: config.eraseZeroWidthCharacter,
    decodeHtmlEntity: config.decodeHtmlEntity,
    regex: config.regex,
  });
  LogInfo(`存储的配置：${JSON.stringify(data)}`);
  console.log("config", data);
  SetOptions({
    eraseZeroWidthCharacter: config.eraseZeroWidthCharacter,
    decodeHtmlEntity: config.decodeHtmlEntity,
    regex: config.regex,
  });
};
</script>

<style lang="less" scoped>
@import "../animation.less";
.menu-area {
  .text-focus-in();
  .menu {
    background-color: transparent;
    box-shadow: 0 0 12px 1px rgb(0 0 0 / 15%);
  }
}
.options-panel {
  transition: 0.25s;
  &::after {
    position: absolute;
    content: "";
    background-image: url(../assets/images/background-image.png);
    background-repeat: no-repeat;
    background-attachment: fixed;
    background-size: cover;
    background-position: top;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 0;
    filter: blur(10px);
    pointer-events: none;
  }
}
.confirm {
  font-size: 1rem;
  font-weight: 700;
  transition: 0.25s;
  &:hover {
    background-color: #fdd9659b;
    border-color: #fdd9659b;
    color: #252525;
  }
}
</style>
