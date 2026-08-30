<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app-root">
    <header class="navbar">
      <h1 style="padding: 10px; color: black; font-family: fantasy">WasaText</h1>
    </header>

    <div class="app-container">
      <RouterView />
    </div>
  </div>
</template>

<style>
.app-root {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header verde */
.navbar {
  background-color: rgb(17, 209, 17);
  margin: 0;
  padding: 0;
  height: 100px;
  display: flex;
  align-items: center;
}

/* Spazio sotto l'header */
.app-container {
  flex: 1;
  display: flex;
  width: 100%;
}
</style>







<!--
<template>
	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<img src="../../public/favicon.ico" alt="logo" class="logo" />
		<span class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-3"> 
			WasaText
		</span>
	</header>

	<div class="container-fluid">
		<RouterView /> 
	</div>
</template>

<style>
</style>
-->