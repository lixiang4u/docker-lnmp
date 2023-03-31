<template>
  <div class="container">
    <Header/>

    <ul class="list">
      <li class="title">
        <span class="p-idx">#</span>
        <span class="p-logs">日志</span>
      </li>
      <li v-for="(item, idx) in logsList" v-bind:key="idx">
        <span class="p-idx">{{ idx + 1 }}</span>
        <span class="p-logs">{{ item }}</span>
      </li>
    </ul>

    <br />
    <br />

  </div>
</template>
<script>
import axios from "axios";
import Header from "@/components/Header.vue";

export default {
  name: "container-logs",
  data() {
    return {
      logsList: [],
      search: {
        containerId: '',
      },
    }
  },
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Header,
  },
  created() {
  },
  mounted() {
    this.search.containerId = this.$route.params['id'] ?? '';
    this.getContainerLogs()
  },
  methods: {
    getContainerLogs() {
      axios.get('/container/logs/' + this.search.containerId).then((response) => {
        if (response.data['code'] === 200) {
          this.logsList = response.data['data']
        }
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
  border-bottom: 1px solid transparent;
  padding: 5px 16px 6px 16px;
  font-family: SourceHanSansSC-regular, "微软雅黑", Arial, Helvetica, sans-serif;
}

.container .list li:hover {
  border-bottom: 1px solid #c3c3c3;
}

.container li > span {
  display: inline-block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0 6px 0 6px;
}

.container li > span.p-idx {
  width: 50px;
}

.container li > span.p-logs {
  width: 980px;
}

</style>