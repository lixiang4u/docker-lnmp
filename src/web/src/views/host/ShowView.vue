<template>
  <div class="container">
    <Header/>

    <div class="form">
      <el-form
          label-width="120px"
          :label-position="'right'"
      >
        <el-form-item label="ID:">
          <el-input v-model="hostInfo.id" disabled placeholder="ID" />
        </el-form-item>
        <el-form-item label="名称：">
          <el-col>
            <el-input v-model="hostInfo.name" placeholder="请填写站点名称，如：默认站点" />
          </el-col>
        </el-form-item>
        <el-form-item label="域名：">
          <el-input v-model="hostInfo.domain" placeholder="请填写站点域名，如：api.localhost.me"/>
        </el-form-item>
        <el-form-item label="项目路径：">
          <el-input v-model="hostInfo.root" placeholder="请填写站点域名，如：/path/to/project/root" @click="onSelectPath(hostInfo.root)" />
        </el-form-item>
        <el-form-item label="web路径：">
          <el-input v-model="hostInfo.web_root" placeholder="请填写站点域名，如：/" :value="hostInfo.web_root??'/'"/>
        </el-form-item>

        <el-form-item>
          <el-button v-if="isCreate" type="primary" @click="onCreate">新增</el-button>
          <el-button v-if="isUpdate" type="primary" @click="onUpdate">修改</el-button>
        </el-form-item>

      </el-form>
    </div>

    <el-dialog v-model="dialogVisible" title="选择目录" draggable>
      <div class="current-root">
        <font-awesome-icon :icon="['fas', 'folder']"/>&nbsp;
        <span v-for="(item,idx) in crumbs" v-bind:key="idx">
          <span class="crumb" @click="onSelectFolder(item['path'],true)">
            {{ item['base'] != '/' ? item['base'] : '' }}
        </span>/</span>
      </div>
      <el-form class="file-list">
        <ul>
          <li>
            <span class="checkbox">#</span>
            <span class="name">文件名</span>
            <span class="modify-time">修改时间</span>
            <span class="permission">权限</span>
          </li>
          <li v-for="(item, idx) in fileList" v-bind:key="idx">
            <span class="checkbox">
              <el-checkbox-group>
                <input type="checkbox" :name="item.path" v-model="item.checked" @click="onFileSelected(item.path)" :disabled="!item.is_dir">
              </el-checkbox-group>
            </span>
            <span class="name">
              <font-awesome-icon :icon="['fas', 'folder']" v-if="item.is_dir"/>
              <font-awesome-icon :icon="['fas', 'file']" v-else/>&nbsp;
              <span @click="onSelectFolder(item.path, item.is_dir)" :class="item.is_dir?'dir':'file'">
                {{ item.name }}
              </span>
            </span>
            <span class="modify-time">{{ formatTime(item.time) }}</span>
            <span class="permission">{{ item.perm }}</span>
          </li>
        </ul>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="dialogVisible = false; hostInfo.root = selectedRoot;">选择</el-button>
      </span>
      </template>
    </el-dialog>

  </div>
</template>
<script>
import axios from "axios";
import {ElMessage} from 'element-plus'
import Header from "@/components/Header.vue";

import {library} from '@fortawesome/fontawesome-svg-core';
import {faFile, faFolder} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';

library.add(faFile, faFolder)

export default {
  name: "host-list",
  data() {
    return {
      root: '',
      crumbs: [],
      fileList: [],
      hostInfo: {},
      isCreate: false,
      isUpdate: false,
      dialogVisible: false,
      selectedRoot: '',
    }
  },
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Header,
    'font-awesome-icon': FontAwesomeIcon,
  },
  created() {
  },
  mounted() {
    //console.log('[this.$route]', this.$route)
    switch (this.$route.query['op']) {
      case 'create':
        this.isCreate = true
        break;
      case 'update':
        this.isUpdate = true
        this.getVirtualHost(this.$route.params['domain'])
        break
      default:
        this.getVirtualHost(this.$route.params['domain'])
    }
  },
  methods: {
    getVirtualHost(domain) {
      axios.get('/host/show/' + domain).then((response) => {
        console.log('[data]', response)
        if (response.data['code'] === 200) {
          this.hostInfo = response.data['data']
        }
      })
    },
    handleFiles(files) {
      console.log('[files]', files)
    },
    onCreate() {
      axios.post('/host/create', this.hostInfo, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }).then((response) => {
        if (response.data['code'] === 200) {
          ElMessage({message: response.data['msg'], type: 'success',})
          setTimeout(() => {
            this.$router.push({name: 'projectList'})
          }, 3000)
        } else {
          ElMessage.error(response.data['msg'])
        }
      })
    },
    onUpdate() {
      axios.put('/host/update/'+this.hostInfo['id'], this.hostInfo, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }).then((response) => {
        console.log('[data]', response)
        if (response.data['code'] === 200) {
          ElMessage({message: response.data['msg'], type: 'success',})
          setTimeout(() => {
            this.$router.push({name: 'projectList'})
          }, 3000)
        } else {
          ElMessage.error(response.data['msg'])
        }
      })
    },
    formatTime(timestamp) {
      return (new Date(timestamp * 1000).toLocaleString());
    },
    onSelectPath(path) {
      console.log('[onSelectPath]')
      this.dialogVisible = true
      this.getFiles(path ?? '')
    },
    getFiles(root) {
      axios.get('/file/list?path=' + root).then((response) => {
        console.log('[data]', response)
        if (response.data['code'] === 200) {
          this.fileList = response.data['data']['files']
          this.root = response.data['data']['root']
          this.crumbs = response.data['data']['crumbs']
        }
      })
    },
    onFileSelected(path) {
      for (const key in this.fileList) {
        if (this.fileList[key].path !== path) {
          this.fileList[key].checked = false
        } else {
          this.fileList[key].checked = true
          console.log('[selected]', path)
          this.selectedRoot = path
        }
      }
    },
    onSelectFolder(path, isDir = false) {
      if (isDir === false) {
        return
      }
      this.getFiles(path)
    },
  }
}
</script>
<style scoped>
.form {
  padding-top: 40px;
}

.current-root {
  margin: -10px 15px 0 15px;
  line-height: 180%;
}

.current-root .crumb {
  color: #337ecc;
  cursor: pointer;
}

.file-list ul {
  padding-left: 0;
  /*margin-top: -10px;*/
}

.file-list li {
  list-style-type: none;
  flex-wrap: nowrap;
  border-bottom: 1px solid #e8e8e8;
  padding: 16px 16px 16px 16px;

}

.file-list li > span {
  display: inline-block;
}

.file-list .checkbox {
  width: 40px;
}

.file-list .name {
  width: 260px;
}

.file-list .name .dir {
  cursor: pointer;
}

.file-list .modify-time {
  width: 180px;
}

.file-list .permission {
  width: 50px;
}

</style>