<template>
  <el-menu
      class="el-menu-demo"
      mode="horizontal"
  >
    <el-menu-item index="1" @click="this.$router.push({name: 'containerList'})">
      <el-text class="mx-1" size="large">docker-lnmp</el-text>
    </el-menu-item>
    <el-menu-item index="2" @click="this.$router.push({name: 'imageList', query: {project: 'all'}})">
      镜像列表
    </el-menu-item>
    <el-menu-item index="3" @click="this.$router.push({name: 'containerList'})">容器列表</el-menu-item>
    <el-menu-item index="4" @click="this.$router.push({name: 'hostList'})">站点列表</el-menu-item>
    <el-menu-item index="5" @click="goCreate">添加站点</el-menu-item>
  </el-menu>
  <br />
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
    goCreate() {
      this.$router.push({name: 'hostCreate', query: {op: 'create'}})
    },
  }
}
</script>

<style scoped>
.mx-1 {
  font-weight: bold;
  font-size: 150%;
}
</style>
