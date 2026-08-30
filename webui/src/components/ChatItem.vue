<template>
  <div class="chat-item" :class="{ active }" @click="$emit('click')">
    <img class="chat-photo" :src="chat.photoUrl || 'https://via.placeholder.com/50'" />

    <div class="chat-info">
      <p class="chat-name">{{ chat.name }}</p>
      <p class="chat-last">{{ chat.lastMessage }}</p>
    </div>

    <div class="chat-time">
      {{ chat.lastMessageTime }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    chat: Object,
    active: Boolean
  }
};
</script>

<style scoped>
.chat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}

.chat-item.active {
  background-color: #f0f0f0;
}

.chat-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
}

.chat-info {
  flex: 1;
}

.chat-name {
  font-weight: bold;
}

.chat-last {
  font-size: 0.9rem;
  color: #666;
}

.chat-time {
  font-size: 0.8rem;
  color: #777;
}
</style>
