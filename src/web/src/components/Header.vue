<template>
  <el-menu
      class="el-menu-demo"
      mode="horizontal"
  >
    <el-menu-item index="1">DOCKER-LNMP</el-menu-item>

    <el-menu-item index="3" @click="goList">站点列表</el-menu-item>
    <el-menu-item index="4" @click="goCreate">添加站点</el-menu-item>
  </el-menu>
</template>

<script>

export default {
  name: "site-header",
  data() {
    return {}
  },
  created() {
  },
  mounted() {
  },
  methods: {
    goList() {
      this.$router.push({name: 'hostList'})
    },
    goCreate() {
      this.$router.push({name: 'hostCreate', query: {op: 'create'}})
    },
  }
}
</script>

<style scoped>
</style>
