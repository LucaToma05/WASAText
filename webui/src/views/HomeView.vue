<template>
  <div class="home-layout">
    <div class="home-sidebar"></div>
    <div class="home-main-area"></div>
  </div>
</template>

<script>
export default {};
</script>

<style scoped>
.home-layout {
  display: flex;
  width: 100%;
  height: 100%;
}

.home-sidebar {
  width: 30%;
  border-right: 2px solid #ccc;
  background-color: white;
}

.home-main-area {
  width: 70%;
  background-color: #d9f7ff;
}
</style>
