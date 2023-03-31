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
              <font-awesome-icon v-if="item['state']==='exited'" :icon="['fas', 'stop']" class="stop"
                                 @click="start(item['id'])"/>
              <font-awesome-icon v-if="item['state']==='running'" :icon="['fas', 'play']" class="start"
                                 @click="stop(item['id'])"/>
            </el-tooltip>
        </span>
        <span>
          <a href="#">
            <el-tooltip placement="top" :content="listToString(item['ports'])">{{
                listToString(item['ports'])
              }}</el-tooltip>
          </a>
        </span>
        <span>
            <el-tooltip placement="top" :content="item['status']">{{ cutLastUpdateTime(item['status']) }}</el-tooltip>
        </span>
        <span>
          [<el-link type="primary" href="#">详情</el-link>]
          [<el-link type="warning" href="#">编辑</el-link>]
          [<el-link type="danger" href="#">删除</el-link>]
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
import {library} from '@fortawesome/fontawesome-svg-core';
import {faPlay, faStop} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';
import {ElLoading, ElMessage} from "element-plus";

library.add(faPlay, faStop)

export default {
  name: "host-list",
  data() {
    return {
      composeStatusList: [],
      dialogVisible: false,
      removeId: null,
      loading: false,
    }
  },
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Header,
    // eslint-disable-next-line vue/no-unused-components
    VideoPlay,
    // eslint-disable-next-line vue/no-unused-components
    VideoPause,
    'font-awesome-icon': FontAwesomeIcon,
  },
  created() {
  },
  mounted() {
    this.getVirtualHost()
  },
  methods: {
    getVirtualHost() {
      axios.get('/container/list').then((response) => {
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
      return value.join(",")
    },
    start(id) {
      this.showLoading()
      axios.post('/container/start/' + id).then((response) => {
        if (response.data['code'] === 200) {
          ElMessage({message: response.data['msg'], type: 'success',})
        } else {
          ElMessage.error(response.data['msg'])
        }
      }).finally(() => {
        this.loading.close()
        this.getVirtualHost()
      })
    },
    stop(id) {
      this.showLoading()
      axios.post('/container/stop/' + id).then((response) => {
        if (response.data['code'] === 200) {
          ElMessage({message: response.data['msg'], type: 'success',})
        } else {
          ElMessage.error(response.data['msg'])
        }
      }).finally(() => {
        this.loading.close()
        this.getVirtualHost()
      })
    },
    showLoading() {
      this.loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
      })
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
  width: 50px;
}

.container li > span:nth-child(2) {
  width: 200px;
}

.container li > span:nth-child(3) {
  width: 200px;
}

.container li > span:nth-child(4) {
  width: 70px;
  text-align: center;
}

.container li > span:nth-child(5) {
  width: 150px;
}

.container li > span:nth-child(6) {
  width: 130px;
}

.stop {
  font-size: 20px;
  color: #6b6969;
  cursor: pointer;
}

.start {
  font-size: 20px;
  color: rgba(0, 128, 0, 0.73);
  cursor: pointer;
}

</style>