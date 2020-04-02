<template>
    <div class="clearfix">
        <h2>管理提醒</h2>

        <div class="col-left">

            <div class="operate clearfix">
                <BatchHandle/>
            </div>

            <table class="list">
                <tbody>
                <tr>
                    <th width="20">&nbsp;</th>
                    <th width="60">提醒类别</th>
                    <th width="180">时间</th>
                    <th>内容</th>
                    <th width="80">操作</th>
                </tr>
                <tr v-if="list.length === 0">
                    <td colspan="7" align="center">还没有提醒！</td>
                </tr>
                <tr v-if="list.length > 0" v-for="(v,k) in list" :key="v.id">
                    <td><input type="checkbox" name="ids" :value="v.id"/></td>
                    <td>{{ remindType[v.type] }}</td>
                    <td>{{ remindTimeFormat(v) }}</td>
                    <td>{{ v.content }}</td>
                    <td><a href="javascript:void(0)" @click="editItem(v.id)">编辑</a><span class="line">|</span><a
                            href="javascript:void(0)"
                            @click="deleteItem(v.id)">删除</a>
                    </td>
                </tr>
                </tbody>
            </table>
            <div class="operate clearfix">
                <BatchHandle/>
                <Paginate
                        v-model="page"
                        :page-count="pageTotal"
                        :click-handler="changePage"
                        :prev-text="'<上一页'"
                        :next-text="'下一页>'"
                        :container-class="'paginator'">
                </Paginate>
            </div>
        </div>
        <div class="col-right" style="width: 250px; padding-top: 31px;">
            <form class="vf" method="post" autocomplete="off" @submit.prevent="submit">
                <p><label class="label_input">提醒类别</label>
                    <select name="type" v-model="item.type">
                        <option :value="k" v-for="(v,k) in remindType" :key="k">{{v}}</option>
                    </select>
                </p>

                <p><label class="label_input">提醒时间</label>
                    <select v-show="[0,6].includes(intRemindType)" v-model="item.month">
                        <option v-for="m in 12" :value="m">{{monthFormat[m]}}月</option>
                    </select>
                    <select v-show="[4].includes(intRemindType)" v-model="item.week">
                        <option v-for="d in 7" :value="d">周{{weekFormat[d]}}</option>
                    </select>
                    <select v-show="[0,5,6].includes(intRemindType)" v-model="item.day">
                        <option v-for="d in 31" :value="d">{{numFormat(d)}}日</option>
                    </select>
                    <select v-show="[0,3,4,5,6].includes(intRemindType)" v-model="item.hour">
                        <option v-for="d in 24" :value="d-1">{{numFormat(d-1)}}时</option>
                    </select>
                    <select v-show="[0,2,3,4,5,6].includes(intRemindType)" v-model="item.minute">
                        <option v-for="d in 60" :value="d-1">{{numFormat(d-1)}}分</option>
                    </select>
                </p>

                <p>
                    <label class="label_input">提醒内容</label>
                    <textarea name="content" rows="5" cols="30" v-model="item.content"></textarea>
                </p>
                <p class="act">
                    <button class="formbutton" type="submit">{{item.id ? '修改':'添加'}}</button>
                    <a v-show="item.id" class="ml10" href="javascript:void(0)" @click="cancel">取消</a>
                </p>
            </form>
        </div>
    </div>
</template>

<script>
  import {remindDeleteApi, remindListApi, remindPostApi} from "../service";
  import Paginate from 'vuejs-paginate'
  import {BatchHandle} from "../components";
  import dayjs from "dayjs"
  import list from "../mixins/list"

  export default {
    name: "AdminRemind",
    data() {
      return {
        listApi: remindListApi,
        postApi: remindPostApi,
        deleteApi: remindDeleteApi,
        item: {},
        defaultRemind() {
          return {
            type: 0,
            month: 1,
            week: 1,
            day: 1,
            hour: 0,
            minute: 0,
          }
        },
        remindType: {
          0: "固定",
          1: "每分钟",
          2: "每小时",
          3: "每天",
          4: "每周",
          5: "每月",
          6: "每年"
        },
        monthFormat: {
          1: '01',
          2: '02',
          3: '03',
          4: '04',
          5: '05',
          6: '06',
          7: '07',
          8: '08',
          9: '09',
          10: '10',
          11: '11',
          12: '12',
        },
        weekFormat: {
          1: '一',
          2: '二',
          3: '三',
          4: '四',
          5: '五',
          6: '六',
          7: '日',
        }
      }
    },
    mixins: [list],
    components: {
      Paginate,
      BatchHandle
    },
    computed: {
      intRemindType() {
        return parseInt(this.item.type)
      }
    },
    methods: {
      numFormat(n) {
        return n < 10 ? '0' + n : n;
      },
      remindTimeFormat(v) {
        let str = ''
        switch (v.type) {
          case 0:
            str = dayjs(v.created_at).year() + "年" + this.monthFormat[v.month] + "月" + this.numFormat(v.day) + "日 " + this.numFormat(v.hour) + "时" + this.numFormat(v.minute) + "分"
            break;
          case 1:
          case 2:
            break;
          case 3:
            str = this.numFormat(v.hour) + "时" + this.numFormat(v.minute) + "分"
            break;
          case 4:
            str = "周" + this.weekFormat[v.week] + " " + this.numFormat(v.hour) + "时" + this.numFormat(v.minute) + "分"
            break;
          case 5:
            str = this.numFormat(v.day) + "日 " + this.numFormat(v.hour) + "时" + this.numFormat(v.minute) + "分"
            break;
          case 6:
            str = this.monthFormat[v.month] + "月" + this.numFormat(v.day) + "日 " + this.numFormat(v.hour) + "时" + this.numFormat(v.minute) + "分"
            break;
        }
        return str
      },
      cancel() {
        this.item = this.defaultRemind()
      },
      submit() {
        let {id, type, content, month, week, day, hour, minute} = this.item
        this.triggerSubmit({id, type: parseInt(type), content, month, week, day, hour, minute}, this.defaultRemind())
      },
    },
    mounted() {
      this.item = this.defaultRemind()
      this.loadList()
    }
  }
</script>

<style scoped>

</style>