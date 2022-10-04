import { createApp } from "vue";
import App from "./App.vue";
import "virtual:svg-icons-register";
import svgIcon from "./components/SvgIcon.vue";
import "./index.css";

createApp(App).component("svg-icon", svgIcon).mount("#app");
