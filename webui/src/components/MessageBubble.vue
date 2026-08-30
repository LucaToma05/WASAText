<template>
  <div class="message" :class="{ mine }">
    {{ msg.content }}
  </div>
</template>

<script>
export default {
  props: {
    msg: Object,
    mine: Boolean
  }
};
</script>

<style scoped>
.message {
  max-width: 60%;
  padding: 10px;
  margin-bottom: 10px;
  background-color: #e6e6e6;
  border-radius: 8px;
}

.message.mine {
  background-color: #cce5ff;
  margin-left: auto;
}
</style>
