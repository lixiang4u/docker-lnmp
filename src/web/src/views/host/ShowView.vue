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
          <el-input v-model="hostInfo.root" placeholder="请填写站点域名，如：/path/to/project/root"/>
        </el-form-item>
        <el-form-item label="web路径：">
          <el-input v-model="hostInfo.web_root" placeholder="请填写站点域名，如：/" />
        </el-form-item>

        <el-form-item>
          <el-button v-if="isCreate" type="primary" @click="onCreate">新增</el-button>
          <el-button v-if="isUpdate" type="primary" @click="onUpdate">修改</el-button>
        </el-form-item>

      </el-form>
    </div>
  </div>
</template>
<script>
import axios from "axios";
import {ElMessage} from 'element-plus'
import Header from "@/components/Header.vue";

export default {
  name: "host-list",
  data() {
    return {
      hostInfo: {},
      isCreate: false,
      isUpdate: false,
    }
  },
  components: {
    Header
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
      console.log('[hostInfo]', this.hostInfo)
      axios.post('/host/create', this.hostInfo, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }).then((response) => {
        console.log('[data]', response)
        if (response.data['code'] === 200) {
          this.hostInfo = response.data['data']
          ElMessage({message: response.data['msg'], type: 'success',})
        } else {
          ElMessage.error(response.data['msg'])
        }
      })
    },
    onUpdate() {
      console.log('[hostInfo]', this.hostInfo)
      axios.put('/host/update/'+this.hostInfo['id'], this.hostInfo, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }).then((response) => {
        console.log('[data]', response)
        if (response.data['code'] === 200) {
          this.hostInfo = response.data['data']
          ElMessage({message: response.data['msg'], type: 'success',})
        } else {
          ElMessage.error(response.data['msg'])
        }
      })
    },
  }
}
</script>
<style scoped>
.form {
  padding-top: 40px;
}

</style>