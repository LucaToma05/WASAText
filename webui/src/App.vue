<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app-root">
    <header class="navbar">
      <h1>WasaText</h1>
    </header>

    <div class="app-container">
      <RouterView />
    </div>
  </div>
</template>

<style>
.app-root {
  height: 100vh;
  height: 100dvh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* Header verde */
.navbar {
  background-color: rgb(17, 209, 17);
  margin: 0;
  padding: 0 20px;
  height: clamp(50px, 8vh, 80px);
  display: flex;
  align-items: center;
  flex-shrink: 0;
  box-sizing: border-box;
}

.navbar h1 {
  padding: 10px;
  color: black;
  font-family: fantasy;
  font-size: clamp(18px, 4vw, 28px);
  margin: 0;
  white-space: nowrap;
}

/* Spazio sotto l'header */
.app-container {
  flex: 1;
  display: flex;
  width: 100%;
  min-height: 0;
  overflow: hidden;
}

@media (max-width: 600px) {
  .navbar {
    height: 48px;
    padding: 0 12px;
  }
}
</style>