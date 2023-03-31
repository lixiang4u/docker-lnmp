<template>
  <div class="container">
    <Header/>

    <ul class="list">
      <li class="title">
        <span class="p-idx">#</span>
        <span class="p-project">项目</span>
        <span class="p-name">容器名称</span>
        <span class="p-image">镜像</span>
        <span class="p-state">状态</span>
        <span class="p-ports">端口</span>
        <span class="p-status">最近启动时间</span>
        <span class="p-op">操作</span>
      </li>
      <div class="project-block" v-for="(tmpList, projectName, projectIdx) in containerList" v-bind:key="projectName">
        <li class="project">
          <span class="p-idx">{{ projectIdx + 1 }}</span>
          <span class="p-project">{{ projectName }}</span>
          <span class="p-name"></span>
          <span class="p-image"></span>
          <span class="p-state">({{ getActiveCount(tmpList) }}/{{ tmpList.length }})</span>
          <span class="p-ports"></span>
          <span class="p-status"></span>
          <span class="p-op">
            <font-awesome-icon :icon="['fas', 'play']" class="start" @click="onProjectStart"/>&nbsp;
            <font-awesome-icon :icon="['fas', 'pause']" class="stop" @click="onProjectStop"/>&nbsp;
            <font-awesome-icon :icon="['fas', 'trash']" class="stop" @click="onProjectRemoveConfirm"/>
          </span>
        </li>
        <li v-for="(item, idx) in tmpList" v-bind:key="idx">
          <span class="p-idx"></span>
          <span class="p-project"></span>
          <span class="p-name">
            <el-tooltip placement="top" :content="item['name']+' : '+item['id']">{{ item['name'] }}</el-tooltip>
          </span>
          <span class="p-image">
            <router-link :to="{name: 'imageList', query: {project: 'all'}}">
              <el-tooltip placement="top" :content="item['image']">{{ item['image'] }}</el-tooltip>
            </router-link>
          </span>
          <span class="p-state">
            <el-tooltip placement="top" :content="item['state']">
              <font-awesome-icon v-if="item['state']==='exited'" :icon="['fas', 'pause']" class="stop"
                                 @click="start(item['id'])"/>
              <font-awesome-icon v-if="item['state']==='running'" :icon="['fas', 'play']" class="start"
                                 @click="stop(item['id'])"/>
            </el-tooltip>
          </span>
          <span class="p-ports">
            <el-tooltip placement="top" :content="portsToText(item['ports'])">{{
                portsToText(item['ports'])
              }}</el-tooltip>
          </span>
          <span class="p-status p-status-text">
            <el-tooltip placement="top" :content="item['status']">{{ (item['status']) }}</el-tooltip>
          </span>
          <span class="p-op">
            [<el-link type="warning" :href="'/container/logs/'+item['id']">日志</el-link>]
            [<el-link type="danger" href="#">删除</el-link>]
          </span>
        </li>
      </div>
    </ul>

    <br/><br/>

    <el-row>
      <el-button type="primary" @click="startAll">一键启动</el-button>
      <el-button type="warning" @click="stopAll">一键停止</el-button>
      <el-button type="info" @click="rebuildConfirm">重新构建</el-button>
    </el-row>

    <el-dialog
        v-model="dialogVisible"
        title="提示"
        width="30%"
    >
      <span>确定重新构建吗？</span>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="onRebuild">
          确认
        </el-button>
      </span>
      </template>
    </el-dialog>

    <br/><br/>

  </div>
</template>
<script>
import axios from "axios";
// import {ElMessage} from "element-plus";
import Header from "@/components/Header.vue";
import {VideoPause, VideoPlay} from "@element-plus/icons-vue";
import {library} from '@fortawesome/fontawesome-svg-core';
import {faPause, faPlay, faStop, faTrash} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';
import {ElLoading, ElMessage} from "element-plus";

library.add(faPlay, faStop, faPause, faTrash)

export default {
  name: "host-list",
  data() {
    return {
      containerList: [],
      dialogVisible: false,
      removeId: null,
      loading: false,
      search: {
        imageId: '',
        project: '',
      },
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
    this.search.imageId = this.$route.query['image_id'] ?? '';
    this.search.project = this.$route.query['project'] ?? '';
    this.getVirtualHost()
  },
  methods: {
    getVirtualHost() {
      const q = {
        project: this.search.project,
        image_id: this.search.imageId,
      }
      axios.get('/project/list', {params: q}).then((response) => {
        if (response.data['code'] === 200) {
          this.containerList = response.data['data']
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
    startAll() {
      for (const key in this.containerList) {
        this.start(this.containerList[key]['id'])
      }
    },
    stopAll() {
      for (const key in this.containerList) {
        this.stop(this.containerList[key]['id'])
      }
    },
    rebuildConfirm() {
      this.dialogVisible = true

    },
    onRebuild() {
      this.dialogVisible = false
      ElMessage({message: '重构中', type: 'success',})
    },
    getActiveCount(containerList) {
      let count = 0
      for (const tmpKey in containerList) {
        if (containerList[tmpKey]['state'] === 'running') {
          count++
        }
      }
      return count
    },
    portsToText(ports) {
      if (!ports) {
        return '-'
      }
      return ports.join(',')
    },
    onProjectStart() {
      //
    },
    onProjectStop() {
      //
    },
    onProjectRemoveConfirm() {
      //
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

.container .list li:hover {
  border-bottom: 1px solid #c3c3c3;
}

.project-block:hover {
  border-left: 6px solid #c3c3c3;
}

.project-block:hover li {
  padding-left: 10px;
}

.container li > span {
  display: inline-block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0 6px 0 6px;
}

.container li > span.p-idx {
  width: 24px;
}

.container li > span.p-project {
  width: 140px;
}

.container li > span.p-id {
  width: 80px;
}

.container li > span.p-name {
  width: 200px;
}

.container li > span.p-image {
  width: 200px;
}

.container li > span.p-state {
  width: 120px;
  text-align: center;
}

.container li > span.p-ports {
  width: 180px;
}

.container li > span.p-status {
  width: 100px;
}

.container li > span.p-status-text {
  font-size: 12px;
  color: #6b6969;
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