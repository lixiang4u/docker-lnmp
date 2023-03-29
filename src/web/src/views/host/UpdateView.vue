<template>
  <div class="container">
    <h1>站点列表</h1>
    <ul>
      <li class="title">
        <span>序号</span>
        <span>名称</span>
        <span>域名</span>
        <span>根目录</span>
        <span>web目录</span>
        <span>操作</span>
      </li>
      <li v-for="(item, idx) in hostList">
        <span>{{ idx + 1 }}</span>
        <span :title="item['name']">{{ item['name'] }}</span>
        <span :title="item['domain']">
          <a :href="'http://'+item['domain']" target="_blank">{{ item['domain'] }}</a>
        </span>
        <span :title="item['root']">
          <a href="#">{{ item['root'] }}</a>
        </span>
        <span :title="item['web_root']">{{ item['web_root'] || '/' }}</span>
        <span>
          [<a href="#">详情</a>]
          [<a href="#">编辑</a>]
          [<a href="#">删除</a>]
        </span>
      </li>
      <li v-for="(item, idx) in hostList">
        <span>{{ idx + 1 }}</span>
        <span>{{ item['name'] }}</span>
        <span>{{ item['domain'] }}</span>
        <span>{{ item['root'] }}</span>
        <span>{{ item['web_root'] || '/' }}</span>
        <span>
          [<a href="#">详情</a>]
          [<a href="#">编辑</a>]
          [<a href="#">删除</a>]
        </span>
      </li>
    </ul>
  </div>
</template>
<script>
import axios from "axios";

export default {
  name: "host-list",
  data() {
    return {
      hostList: []
    }
  },
  created() {
    console.log('AAAA')
  },
  mounted() {
    this.getVirtualHost()
  },
  methods: {
    getVirtualHost() {
      axios.get('/host/list').then((response) => {
        console.log('[data]', response)
        if (response.data['code'] === 200) {
          this.hostList = response.data['data']
        }
      })
    }
  }
}
</script>
<style scoped>
.container {
  width: 1024px;
  font-family: "Source Code Pro", "微软雅黑", Arial, Helvetica, sans-serif;
}

.container ul {
  padding-left: 0;
}

.container li.title span {
  font-weight: bold;
  /*text-align: center;*/
}

.container li {
  list-style-type: none;
  flex-wrap: nowrap;
  border-bottom: 1px solid #c3c3c3;
  padding: 0 16px 0 16px;
  font-family: SourceHanSansSC-regular, "微软雅黑", Arial, Helvetica, sans-serif;
}

.container li span {
  line-height: 320%;
  display: inline-block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0 5px 0 5px;
}

.container li span:nth-child(1) {
  width: 50px;
}

.container li span:nth-child(2) {
  width: 120px;
}

.container li span:nth-child(3) {
  width: 180px;
}

.container li span:nth-child(4) {
  width: 320px;
}

.container li span:nth-child(5) {
  width: 120px;
}

</style>