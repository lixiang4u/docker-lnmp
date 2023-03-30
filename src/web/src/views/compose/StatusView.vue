<template>
  <div class="container">
    <Header/>


    <ul class="list">
      <li class="title">
        <span>序号</span>
        <span>名称</span>
        <span>镜像</span>
        <span>状态</span>
        <span>端口</span>
        <span>最近启动时间</span>
        <span>操作</span>
      </li>
      <li v-for="(item, idx) in composeStatusList" v-bind:key="idx">
        <span>{{ idx + 1 }}</span>
        <span>
          <el-tooltip placement="top" :content="item['name']+' : '+ item['id']">{{ item['name'] }}</el-tooltip>
        </span>
        <span>
          <a href="#">
            <el-tooltip placement="top" :content="item['image']">{{ item['image'] }}</el-tooltip>
          </a>
        </span>
        <span>
            <el-tooltip placement="top" :content="item['state']">
              <el-icon><VideoPlay/></el-icon>
              <el-icon><VideoPause/></el-icon>
              {{ item['state'] }}
            </el-tooltip>
        </span>
        <span>
          <a href="#">
            <el-tooltip placement="top" :content="listToString(item['ports'])">{{ listToString(item['ports']) }}</el-tooltip>
          </a>
        </span>
        <span>
            <el-tooltip placement="top" :content="item['status']">{{ cutLastUpdateTime(item['status']) }}</el-tooltip>
        </span>
        <span>
          [<el-link type="primary" :href="'/host/show/'+item['id']+'?op=show'">详情</el-link>]
          [<el-link type="warning" :href="'/host/update/'+item['id']+'?op=update'">编辑</el-link>]
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
// import {ElMessage} from "element-plus";
import Header from "@/components/Header.vue";
import {VideoPause, VideoPlay} from "@element-plus/icons-vue";

export default {
  name: "host-list",
  data() {
    return {
      composeStatusList: [],
      dialogVisible: false,
      removeId: null,
    }
  },
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Header,
    VideoPlay,
    VideoPause,
  },
  created() {
  },
  mounted() {
    this.getVirtualHost()
  },
  methods: {
    getVirtualHost() {
      axios.get('/compose/status').then((response) => {
        if (response.data['code'] === 200) {
          this.composeStatusList = response.data['data']
        }
      })
    },
    cutLastUpdateTime(value) {
      let arr = value.split(") ")
      if (arr.length > 1) {
        return arr[1]
      }
      return arr[0]
    },
    listToString(value) {
      if (value == null) {
        return ''
      }
      return value.join("\r\n")
    }
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
  width: 50px;
}

.container li > span:nth-child(2) {
  width: 200px;
}

.container li > span:nth-child(3) {
  width: 200px;
}

.container li > span:nth-child(4) {
  width: 100px;
}

.container li > span:nth-child(5) {
  width: 130px;
}

.container li > span:nth-child(6) {
  width: 120px;
}

</style>