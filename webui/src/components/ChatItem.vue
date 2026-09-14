<template>
  <!-- Applica la classe se la chat è selezionata -->
   <!-- Se la chat viene cliccata viene emesso un evento click al componente padre-->
  <div 
    class="chat-preview"
    :class="{ selected: isSelected }"
    @click="$emit('click')"
  >
    <!-- Foto della chat -->
    <img :src="chat.photo || '/default.png'" class="chat-photo" alt="Avatar" />

    <div class="chat-info">
      <div class="chat-top-row">
        <!-- Nome della chat e orario ultimo messaggio -->
        <span class="chat-name">{{ chat.name }}</span>
        <span class="chat-time">{{ formattedTime }}</span>
      </div>

      <div class="chat-bottom-row">
        <div class="chat-last-message">
          <!-- Se l'ultimo messaggio è una foto -->
          <span v-if="parsedContent.isPhoto" class="chat-last-photo">
            <span class="photo-icon">📷</span> Foto
          </span>
          <!-- Se l'ultimo messaggio è un messaggio eliminato -->
          <span v-else-if="parsedContent.isDeleted" class="chat-last-deleted">
            <em>Questo messaggio è stato eliminato</em>
          </span>
          <!-- Altrimenti mostra il testo del messaggio -->
          <span v-else>
            {{ parsedContent.text }}
          </span>
        </div>
        <!-- Indicatore di messaggio non letto -->
        <span v-if="isUnread" class="unread-dot" title="Nuovo messaggio non letto"></span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  // Nome del componente
  name: "ChatItem",
  // Definisce gli oggetti passati
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
  // Dichiara che il componente può emettere l'evento click
  emits: ["click"],

  // Definisce proprietà elaborate al variare degli oggetti passati
  computed: {
    // Formatta l'orario dell'ultimo messaggio in formato HH:MM
    formattedTime() {
      if (!this.chat.lastMessageTime) return "";
      const date = new Date(this.chat.lastMessageTime);
      return isNaN(date.getTime()) ? "" : date.toLocaleTimeString("it-IT", {
        hour: "2-digit",
        minute: "2-digit"
      });
    },
    // Analizza il contenuto dell'ultimo messaggio per determinare se è una foto, un messaggio eliminato o un messaggio normale
    parsedContent() {
      let text = (this.chat.lastMessage || "").toString().trim();

      const isDeleted = text === "Questo messaggio è stato eliminato";
      const isPhoto = !isDeleted && this.checkPhoto(text);

      return { text, isPhoto, isDeleted };
    }
  },
  // Funzioni eseguibili
  methods: {
    // Controlla se il messaggio è una foto usando una regex per identificare le estensioni comuni delle immagini
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
  border-radius: 50%;
  background-color: #25D366;
  flex-shrink: 0;
  margin-left: 8px;
}
</style>
