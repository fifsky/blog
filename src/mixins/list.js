import {dialog, sync} from "../utils";

const list = {
  data () {
    return {
      item:{},
      list: [],
      pageTotal: 0,
      page: 1
    }
  },
  methods: {
    triggerSubmit(data,defaultItem){
      defaultItem = defaultItem || {}
      sync(async () => {
        let ret = await this.postApi(data)
        this.item = defaultItem
        this.$message.success("发表成功")
        this.loadList()
      })
    },
    cancel(){
      this.item = {}
    },
    changePage(pageNum) {
      this.page = pageNum
      let q = {...this.$route.query}
      q.page = pageNum
      this.$router.push({path:this.$router.path,query:q})
      this.loadList()
    },
    editItem(id) {
      this.item = this.list.filter(item=>item.id === id)[0]
    },
    deleteItem(id) {
      dialog.confirm("确认要删除？", (ok) => {
        if (ok) {
          sync(async () => {
            await this.deleteApi({id})
            this.loadList()
          })
        }
      })
    },
    loadList() {
      if (this.$route.query.page) {
        this.page = parseInt(this.$route.query.page)
      }
      let data = {...this.$route.query}
      data.page = this.page
      sync(async () => {
        let ret = await this.listApi(data)
        this.list = ret.list
        this.pageTotal = ret.pageTotal
      })
    },
  }
}
export default list