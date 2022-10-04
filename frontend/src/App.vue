<script setup>
import { onMounted } from "vue";
import Upload from "./components/Upload.vue";
import Options from "./components/Options.vue";

const disabledKeyboardShortcuts = () => {
  window.onkeydown =
    window.onkeyup =
    window.onkeypress =
      (event) => {
        //禁止ctrl+u
        if (event.ctrlKey && event.key.toLowerCase() === "u") return false;
        //禁止ctrl+s
        if (event.ctrlKey && event.key.toLowerCase() === "s") return false;
        //禁止ctrl+shift+i
        if (event.ctrlKey && event.shiftKey && event.key.toLowerCase() === "i")
          return false;
        //禁止 F5
        if (event.key === "F5") return false;
        //禁止 F12
        if (event.key === "F12") return false;
      };
};
const disableContextmenu = () => {
  // 禁用右键
  document.body.onselectstart = document.body.oncontextmenu = () => false;
};

onMounted(() => {
  disabledKeyboardShortcuts();
  disableContextmenu();
});
</script>

<template>
  <div class="main hero relative min-h-screen" data-wails-drag>
    <Upload />
    <Options />
  </div>
</template>

<style lang="less" scoped>
.main {
  background-image: url(./assets/images/background-image.png);
}
</style>
