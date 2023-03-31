<template>
  <div class="container">
    <Header/>

    <ul class="list">
      <li class="title">
        <span class="p-idx">#</span>
        <span class="p-id">ID</span>
        <span class="p-name">镜像名称</span>
        <span class="p-image">TAG</span>
        <span class="p-created">创建时间</span>
        <span class="p-size">大小</span>
        <span class="p-op">操作</span>
      </li>
      <li v-for="(item, idx) in imageList" v-bind:key="idx">
        <span class="p-idx">{{ idx + 1 }}</span>
        <span class="p-id">
          <el-tooltip placement="top" :content="item['id']">{{ item['id'] }}</el-tooltip>
        </span>
        <span class="p-name">
          <a href="#">
            <el-tooltip placement="top" :content="item['tag']">{{ splitTag(item['tag'], 0) }}</el-tooltip>
          </a>
        </span>
        <span class="p-image">
          <el-tooltip placement="top" :content="splitTag(item['tag'],1)">{{ splitTag(item['tag'], 1) }}</el-tooltip>
        </span>
        <span class="p-created">{{ formatTime(item['created_at']) }}</span>
        <span class="p-size p-size-text">{{ formatSize(item['size']) }}</span>
        <span class="p-op">
          [<el-link type="primary" :href="'/container/list/?image_id='+item['id']">详情</el-link>]
          [<el-link type="warning" href="#" style="text-decoration: line-through">编辑</el-link>]
          [<el-link type="danger" @click="removeImageConfirm(item['id'])">删除</el-link>]
        </span>
      </li>
    </ul>

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
      <span>{{ dialogContent }}</span>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="dialogOnClick">
          确认
        </el-button>
      </span>
      </template>
    </el-dialog>

    <br /><br />
  </div>
</template>
<script>
import axios from "axios";
// import {ElMessage} from "element-plus";
import Header from "@/components/Header.vue";
import {VideoPause, VideoPlay} from "@element-plus/icons-vue";
import {library} from '@fortawesome/fontawesome-svg-core';
import {faPause, faPlay, faStop} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';
import {ElLoading, ElMessage} from "element-plus";

library.add(faPlay, faStop, faPause)

export default {
  name: "host-list",
  data() {
    return {
      imageList: [],
      dialogVisible: false,
      dialogContent: '',
      dialogOnClick: '',
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
    // eslint-disable-next-line vue/no-unused-components
    'font-awesome-icon': FontAwesomeIcon,
  },
  created() {
  },
  mounted() {
    this.getVirtualHost()
  },
  methods: {
    getVirtualHost() {
      axios.get('/image/list').then((response) => {
        if (response.data['code'] === 200) {
          this.imageList = response.data['data']
        }
      })
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
      for (const key in this.imageList) {
        this.start(this.imageList[key]['id'])
      }
    },
    stopAll() {
      for (const key in this.imageList) {
        this.stop(this.imageList[key]['id'])
      }
    },
    rebuildConfirm() {
      this.dialogContent = '确定重新构建吗？'
      this.dialogVisible = true
      this.dialogOnClick = this.onRebuild
    },
    onRebuild() {
      this.dialogVisible = false
      ElMessage({message: '重构中', type: 'success',})
    },
    splitTag(tagName, idx) {
      const tmpArr = tagName.split(':')
      if (tmpArr.length > 1) {
        return tmpArr[idx]
      } else {
        return tmpArr[0]
      }
    },
    formatTime(timestamp) {
      return (new Date(timestamp * 1000).toLocaleString()).substring(0, 15);
    },
    formatSize(size, pointLength, units) {
      let unit;
      units = units || ['B', 'K', 'M', 'G', 'TB'];
      while ((unit = units.shift()) && size > 1024) {
        size = size / 1024;
      }
      return (unit === 'B' ? size : size.toFixed(pointLength === undefined ? 2 : pointLength)) + unit;
    },
    removeImageConfirm(id) {
      this.removeId = id
      this.dialogContent = '确定删除镜像吗？'
      this.dialogVisible = true
      this.dialogOnClick = this.onRemoveImage
    },
    onRemoveImage() {
      this.dialogVisible = false
      this.showLoading()
      axios.delete('/image/remove/' + this.removeId).then((response) => {
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
  },
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

.container li > span.p-idx {
  width: 24px;
}

.container li > span.p-id {
  width: 100px;
}

.container li > span.p-name {
  width: 280px;
}

.container li > span.p-image {
  width: 180px;
}

.container li > span.p-created {
  width: 200px;
}

.container li > span.p-size {
  width: 110px;
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