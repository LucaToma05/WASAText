<template>
  <div class="chat-item" :class="{ active, 'has-unread': isUnread }" @click="$emit('click')">
    <img class="chat-photo" :src="chat.photoUrl || chat.photo || 'https://via.placeholder.com/50'" alt="Avatar" />

    <div class="chat-info">
      <p class="chat-name">{{ chat.name }}</p>
      <p class="chat-last">
        <span v-if="isDeletedLastMessage" class="chat-last-deleted"><em>Questo messaggio è stato eliminato</em></span>
        <span v-else-if="isPhotoLastMessage" class="chat-last-photo">📷 Foto</span>
        <span v-else>{{ displayLastMessage }}</span>
      </p>
    </div>

    <div class="chat-meta">
      <div class="chat-time">
        {{ chat.lastMessageTime }}
      </div>
      <span v-if="isUnread" class="unread-dot" title="Nuovo messaggio non letto"></span>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    chat: {
      type: Object,
      required: true
    },
    active: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    isUnread() {
      if (!this.chat || this.active) return false;
      const myId = localStorage.getItem("userId") || localStorage.getItem("token");
      if (this.chat.lastMessageUserId && myId && String(this.chat.lastMessageUserId) === String(myId)) {
        return false;
      }
      return Boolean(this.chat.hasUnread || (this.chat.unreadCount && this.chat.unreadCount > 0));
    },
    isDeletedLastMessage() {
      if (!this.chat || !this.chat.lastMessage) return false;
      return this.chat.lastMessage.toString().trim() === "Questo messaggio è stato eliminato";
    },
    isPhotoLastMessage() {
      if (!this.chat || !this.chat.lastMessage) return false;
      const text = this.chat.lastMessage.toString().trim();
      if (text === "Questo messaggio è stato eliminato") return false;
      if (text.startsWith("📷 Foto")) return true;
      if (text.startsWith("{") && text.endsWith("}")) {
        try {
          const parsed = JSON.parse(text);
          if (parsed.type === "photo" || parsed.isPhoto || parsed.photo) return true;
          if (parsed.text && this.checkPhotoPattern(parsed.text)) return true;
        } catch (e) {}
      }
      return this.checkPhotoPattern(text);
    },
    displayLastMessage() {
      if (!this.chat || !this.chat.lastMessage) return "";
      return this.chat.lastMessage;
    }
  },
  methods: {
    checkPhotoPattern(str) {
      if (!str || typeof str !== "string") return false;
      const t = str.trim();
      return (
        /^(https?:\/\/.*\.(?:png|jpg|jpeg|gif|webp|svg|bmp|avif)(\?.*)?)$/i.test(t) ||
        /^data:image\//i.test(t)
      );
    }
  }
};
</script>

<style scoped>
.chat-item {
  display: flex;
  align-items: center;
  padding: 10px;
  cursor: pointer;
  border-bottom: 1px solid #eee;
  background-color: #ffffff;
  transition: background-color 0.2s ease;
}

.chat-last-deleted em {
  font-style: italic;
  color: #94a3b8;
}

.chat-item:hover {
  background-color: #f8fafc;
}

.chat-item.active {
  background-color: #f0f0f0;
}

.chat-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.chat-info {
  flex: 1;
  min-width: 0;
}

.chat-name {
  font-weight: bold; /* Sempre in grassetto */
  font-size: 1rem;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin: 0;
}

.chat-last {
  font-size: 0.9rem;
  color: #64748b;
  font-weight: 400;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin: 4px 0 0 0;
}

.chat-last-photo {
  color: #475569;
  font-weight: 500;
}

.chat-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
  flex-shrink: 0;
}

.chat-time {
  font-size: 0.8rem;
  color: #777;
  white-space: nowrap;
}

/* Indicatore verde per messaggi non letti */
.unread-dot {
  width: 10px;
  height: 10px;
  min-width: 10px;
  min-height: 10px;
  border-radius: 50%;
  background-color: #25D366;
  flex-shrink: 0;
  display: inline-block;
}
</style>
