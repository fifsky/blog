import Vue from 'vue'
import App from './App.vue'
import './assets/css/base.css'
import router from "./router";
import dayjs from "dayjs"
import isBetween from 'dayjs/plugin/isBetween'
import 'highlight.js/styles/atelier-lakeside-light.css'
import store from "./store"
import Message from 'vue-m-message'
import Meta from 'vue-meta'

import hljs from 'highlight.js/lib/highlight'
import javascript from 'highlight.js/lib/languages/javascript'
import php from 'highlight.js/lib/languages/php'
import go from 'highlight.js/lib/languages/go'
import python from 'highlight.js/lib/languages/python'
import nginx from 'highlight.js/lib/languages/nginx'
import sql from 'highlight.js/lib/languages/sql'
import lua from 'highlight.js/lib/languages/lua'
import bash from 'highlight.js/lib/languages/bash'
import css from 'highlight.js/lib/languages/css'
import java from 'highlight.js/lib/languages/java'
import xml from 'highlight.js/lib/languages/xml'

hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('php', php)
hljs.registerLanguage('go', go)
hljs.registerLanguage('python', python)
hljs.registerLanguage('nginx', nginx)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('lua', lua)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('css', css)
hljs.registerLanguage('java', java)
hljs.registerLanguage('xml', xml)


Vue.use(Message)
Vue.use(Meta)

Vue.config.productionTip = false;

Vue.filter("formatDate", function (v,f) {
    return dayjs(v).format(f)
})

Vue.filter("humanTime",function (v) {
    dayjs.extend(isBetween)
    let currTime = dayjs().add(1,"second")
    let itemTime = dayjs(v)

    if(itemTime.isBetween(currTime.subtract(60,'second'),currTime)){
      return currTime.diff(itemTime,'second')+'秒前'
    }else if(itemTime.isBetween( currTime.subtract(60,'minute'),currTime.subtract(1,'minute'))){
      return currTime.diff(itemTime,'minute')+'分钟前'
    }else if(itemTime.isBetween(currTime.startOf('day'),currTime.endOf('day'))){
      return '今天'+itemTime.format("HH:mm")
    }else if(itemTime.isBetween(currTime.subtract(1,'day').startOf('day'),currTime.subtract(1,'day').endOf('day'))){
      return '昨天'+itemTime.format("HH:mm")
    }else if(itemTime.isBetween(currTime.startOf("year"),currTime.subtract(1,'day').endOf('day'))){
      return itemTime.format("MM月DD日 HH:mm")
    }else{
      return itemTime.format("YYYY-MM-DD HH:mm")
    }
})

Vue.directive('highlight',function (el) {
    let blocks = el.querySelectorAll('pre code');
    blocks.forEach((block)=>{
        hljs.highlightBlock(block)
    })
})

new Vue({
    router,
    store,
    render: h => h(App)
}).$mount('#app');