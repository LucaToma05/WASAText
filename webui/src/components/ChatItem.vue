<template>
  <div 
    class="chat-preview"
    :class="{ 
      selected: isSelected,
      'unread-chat': isUnread
    }"
    @click="$emit('click')"
  >
    <img :src="chat.photo || '/default.png'" class="chat-photo" alt="Avatar" />

    <div class="chat-info">
      <div class="chat-top-row">
        <span class="chat-name">{{ chat.name }}</span>
        <span class="chat-time">{{ formattedTime }}</span>
      </div>

      <div class="chat-bottom-row">
        <div class="chat-last-message">
          <span v-if="isPhoto" class="chat-last-photo">
            <span class="photo-icon">📷</span> Foto
          </span>
          <span v-else-if="isDeleted" class="chat-last-deleted">
            <em>Questo messaggio è stato eliminato</em>
          </span>
          <span v-else>
            {{ displayMessage }}
          </span>
        </div>
        <span v-if="isUnread" class="unread-dot" title="Nuovo messaggio non letto"></span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ChatItem",
  props: {
    chat: {
      type: Object,
      required: true
    },
    isSelected: {
      type: Boolean,
      default: false
    },
    isUnread: {
      type: Boolean,
      default: false
    }
  },
  emits: ["click"],
  computed: {
    formattedTime() {
      if (!this.chat.lastMessageTime) return "";
      let str = this.chat.lastMessageTime.toString().trim();
      if (!str) return "";

      if (!str.includes("T") && str.includes(" ")) {
        str = str.replace(" ", "T") + "Z";
      } else if (str.includes("T") && !str.endsWith("Z") && !str.includes("+") && !str.includes("-", 10)) {
        str += "Z";
      }

      const date = new Date(str);
      if (isNaN(date.getTime())) return "";

      return date.toLocaleTimeString("it-IT", {
        hour: "2-digit",
        minute: "2-digit"
      });
    },
    parsedContent() {
      const text = this.chat.lastMessage;
      if (!text) return { text: "", isPhoto: false, isDeleted: false };
      const raw = String(text).trim();

      if (raw === "Questo messaggio è stato eliminato") {
        return { text: raw, isPhoto: false, isDeleted: true };
      }

      let inner = raw;
      if (raw.startsWith("{") && raw.endsWith("}")) {
        try {
          const parsed = JSON.parse(raw);
          inner = parsed.text || parsed.url || "";
        } catch (e) {
          inner = raw;
        }
      }

      if (inner === "Questo messaggio è stato eliminato") {
        return { text: inner, isPhoto: false, isDeleted: true };
      }

      const isPhoto = this.checkPhoto(inner);
      return { text: inner, isPhoto, isDeleted: false };
    },
    isPhoto() {
      return this.parsedContent.isPhoto;
    },
    isDeleted() {
      return this.parsedContent.isDeleted;
    },
    displayMessage() {
      return this.parsedContent.text;
    }
  },
  methods: {
    checkPhoto(str) {
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
.chat-preview {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 18px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  background-color: #ffffff;
  transition: background-color 0.15s ease, border-left-color 0.15s ease;
  border-left: 4px solid transparent;
}

.chat-preview:hover {
  background-color: #f7f9fa;
}

.chat-preview.selected {
  background-color: #e8f5e9;
  border-left-color: #2e7d32;
}

.chat-photo {
  width: 48px;
  height: 48px;
  min-width: 48px;
  min-height: 48px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  border: 1px solid #e0e0e0;
}

.chat-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.chat-top-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.chat-name {
  font-weight: 600;
  font-size: 15px;
  color: #111;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-time {
  font-size: 12px;
  color: #888;
  flex-shrink: 0;
  margin-left: 8px;
}

.chat-bottom-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-last-message {
  font-size: 13px;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.chat-last-photo {
  color: #444;
  font-style: italic;
}

.chat-last-deleted {
  color: #999;
}

.photo-icon {
  margin-right: 2px;
}

.unread-dot {
  width: 10px;
  height: 10px;
  min-width: 10px;
  min-height: 10px;
  border-radius: 50%;
  background-color: #25D366;
  flex-shrink: 0;
  margin-left: 8px;
  display: inline-block;
}
</style>
