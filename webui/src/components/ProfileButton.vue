<script setup>
import { useRouter } from "vue-router";

const router = useRouter();

function goToProfile() {
    router.push("/profile");
}
</script>

<template>
  <button class="profile-btn" @click="goToProfile">
    Profilo
  </button>
</template>

<style scoped>
.profile-btn {
  padding: 8px 14px;
  cursor: pointer;
  background-color: #ddd;
  border: 1px solid #bbb;
}
</style>
