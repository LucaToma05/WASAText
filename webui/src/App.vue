<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header class="navbar">
    <span>
      <h1 style="padding: 10px; color: black; font-family: fantasy">WasaText</h1>
    </span>
  </header>

  <div class="container" style="overflow: hidden;">
	  <RouterView />
  </div>
</template>

<style>
  .navbar {
    background-color: rgb(17, 209, 17);
    margin: 0;
    padding: 0;
    box-sizing: border-box;
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