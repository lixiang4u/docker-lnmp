<template>
  <div class="container">
    <Header/>

    <ul class="list">
      <li class="title">
        <span>#</span>
        <span>名称</span>
        <span>域名</span>
        <span>根目录</span>
        <span>web目录</span>
        <span>操作</span>
      </li>
      <li v-for="(item, idx) in hostList" v-bind:key="idx">
        <span>{{ idx + 1 }}</span>
        <span>
          <el-tooltip placement="top" :content="item['name']">{{ item['name'] }}</el-tooltip>
        </span>
        <span>
          <a :href="'http://'+item['domain']" target="_blank">
            <el-tooltip placement="top" :content="item['domain']">{{ item['domain'] }}</el-tooltip>
          </a>
        </span>
        <span>
          <el-tooltip placement="top" :content="item['root']">{{ item['root'] }}</el-tooltip>
        </span>
        <span>
          <el-tooltip placement="top" :content="item['web_root'] || '/' ">
            {{ item['web_root'] || '/' }}
          </el-tooltip>
        </span>
        <span>
          [<el-link type="primary" :href="'/host/show/'+item['id']+'?op=show'" target="_blank">详情</el-link>]
          [<el-link type="warning" :href="'/host/update/'+item['id']+'?op=update'" target="_blank">编辑</el-link>]
          [<el-link type="danger" @click="onRemoveConfirm(item['id'])">删除</el-link>]
        </span>
      </li>
    </ul>

    <el-dialog
        v-model="dialogVisible"
        title="提示"
        width="30%"
    >
      <span>确定删除该站点？</span>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="onRemove">
          确认
        </el-button>
      </span>
      </template>
    </el-dialog>
  </div>
</template>
<script>
import axios from "axios";
import {ElMessage} from "element-plus";
import Header from "@/components/Header.vue";

export default {
  name: "host-list",
  data() {
    return {
      hostList: [],
      dialogVisible: false,
      removeId: null,
    }
  },
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Header
  },
  created() {
  },
  mounted() {
    this.getVirtualHost()
  },
  methods: {
    getVirtualHost() {
      axios.get('/host/list').then((response) => {
        if (response.data['code'] === 200) {
          this.hostList = response.data['data']
        }
      })
    },
    onRemoveConfirm(id) {
      this.removeId = id
      this.dialogVisible = true
    },
    onRemove() {
      console.log('[remove]', this.removeId)
      axios.delete('/host/delete/' + this.removeId).then((response) => {
        if (response.data['code'] === 200) {
          ElMessage({message: response.data['msg'], type: 'success',})
        } else {
          ElMessage.error(response.data['msg'])
        }
      }).finally(() => {
        this.dialogVisible = false
        this.removeId = ''
        this.getVirtualHost()
      })
    },
    goList() {
      this.$router.push({name: 'hostShow'})
    },
    goCreate() {
      this.$router.push({name: 'hostCreate', query: {op: 'create'}})
    },
  }
}
</script>
<style scoped>
.container {
}

.container .list {
  padding-left: 0;
}

.container li.title span {
  font-weight: bold;
  /*text-align: center;*/
}

.container .list li {
  list-style-type: none;
  flex-wrap: nowrap;
  border-bottom: 1px solid #e8e8e8;
  padding: 16px 16px 16px 16px;
  font-family: SourceHanSansSC-regular, "微软雅黑", Arial, Helvetica, sans-serif;
}

.container li > span {
  display: inline-block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0 6px 0 6px;
}

.container li > span:nth-child(1) {
  width: 24px;
}

.container li > span:nth-child(2) {
  width: 140px;
}

.container li > span:nth-child(3) {
  width: 200px;
}

.container li > span:nth-child(4) {
  width: 440px;
}

.container li > span:nth-child(5) {
  width: 100px;
}

</style>