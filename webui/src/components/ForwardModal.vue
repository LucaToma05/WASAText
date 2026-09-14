<template>
  <!-- Se cliccato chiude il pannello -->
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-forward-card">
      <div class="modal-card-header">
        <h2 class="modal-title">Inoltra messaggio</h2>
        <button type="button" class="btn-close-modal" @click="$emit('close')">✕</button>
      </div>

      <!-- Anteprima messaggio da inoltrare -->
      <div v-if="message" class="forward-msg-preview-card">
        <div class="forward-msg-preview-label">Contenuto da inoltrare:</div>
        <div class="forward-msg-preview-content">
          <span v-if="isPhoto" class="photo-indicator">
            📷 Foto
          </span>
          <span v-else>{{ displaySnippet }}</span>
        </div>
      </div>

      <!-- Barra di ricerca chat per inoltro -->
      <input 
        v-model="searchQuery"
        type="text"
        placeholder="Cerca conversazione o gruppo a cui inoltrare..."
        class="modal-input"
      />

      <!-- Messaggi di errore o successo -->
      <p v-if="successMsg" class="success-msg">{{ successMsg }}</p>
      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

      <!-- Lista conversazioni destinatarie -->
      <div class="forward-chats-list">
        <div 
          v-for="chat in forwardableChats" 
          :key="chat.chatId"
          class="forward-chat-row"
        >
          <!-- Foto e nome chat -->
          <img :src="chat.photo || '/default.png'" class="forward-chat-avatar" alt="Foto chat" />
          <div class="forward-chat-meta">
            <div class="forward-chat-title">{{ chat.name }}</div>
            <span v-if="chat.isGroup" class="chat-header-badge">Gruppo</span>
          </div>
          <button 
            type="button" 
            class="btn-action-small btn-primary-small"
            :disabled="isForwarding"
            @click="forwardToChat(chat)"
          >
            {{ isForwarding && targetChatId === chat.chatId ? 'Invio...' : 'Inoltra' }}
          </button>
        </div>

        <div v-if="forwardableChats.length === 0" class="no-chats-found">
          Nessun'altra conversazione disponibile
        </div>
      </div>

      <div class="modal-buttons">
        <button type="button" class="btn-secondary" @click="$emit('close')">Chiudi</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ForwardModal",
  props: {
    message: {
      type: Object,
      required: true
    },
    chats: {
      type: Array,
      default: () => []
    },
    currentChatId: {
      type: [String, Number],
      default: null
    }
  },
  emits: ["close", "forwarded"],
  data() {
    return {
      searchQuery: "",
      errorMsg: "",
      successMsg: "",
      isForwarding: false,
      targetChatId: null
    };
  },
  computed: {
    // Filtra le chat a cui inoltrare il messaggio
    forwardableChats() {
      const q = this.searchQuery.trim().toLowerCase();
      return this.chats.filter(c => {
        const isNotCurrent = c.chatId !== this.currentChatId;
        const matchesQuery = !q || (c.name && c.name.toLowerCase().includes(q));
        return isNotCurrent && matchesQuery;
      });
    },
    // Controlla se il messaggio è una foto
    isPhoto() {
      const text = this.message?.text;
      if (!text || typeof text !== "string") return false;
      const t = text.trim();
      return (
        /^(https?:\/\/.*\.(?:png|jpg|jpeg|gif|webp|svg|bmp|avif)(\?.*)?)$/i.test(t) ||
        /^data:image\//i.test(t)
      );
    },
    // Mostra una anteprima del testo del messaggio (tagliandola)
    displaySnippet() {
      const text = this.message?.text;
      if (!text) return "";
      const cleaned = text.replace(/\s+/g, " ").trim();
      return cleaned.length <= 50 ? cleaned : cleaned.slice(0, 50) + "...";
    }
  },
  methods: {
    async forwardToChat(targetChat) {
      if (!this.message || !targetChat) return;
      this.isForwarding = true;
      this.targetChatId = targetChat.chatId;
      this.errorMsg = "";
      this.successMsg = "";

      try {
        await this.$axios.post(
          `/chats/${this.currentChatId}/messages/${this.message.id}/forward`,
          { chatId: targetChat.chatId },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.successMsg = `Messaggio inoltrato con successo a "${targetChat.name}"!`;
        this.$emit("forwarded", targetChat);

        // Chiude il modale dopo 1000ms in modo da far leggere la conferma dell'inoltro del messaggio
        setTimeout(() => {
          this.$emit("close");
        }, 1000);
      } catch (err) {
        this.errorMsg = err.response?.data || "Errore durante l'inoltro del messaggio";
      } finally {
        // Resetta gli stati
        this.isForwarding = false;
        this.targetChatId = null;
      }
    }
  }
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.45);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 500;
}

.modal-forward-card {
  background-color: #ffffff;
  padding: 26px 30px;
  border-radius: 22px;
  width: 480px;
  max-width: 95%;
  max-height: 85vh;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.18);
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-sizing: border-box;
}

.modal-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  font-size: 24px;
  color: #222;
  margin: 0;
}

.btn-close-modal {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #666;
  padding: 4px;
  transition: color 0.15s ease;
}

.btn-close-modal:hover {
  color: #111;
}

.forward-msg-preview-card {
  background-color: #f5f5f5;
  border-left: 4px solid #222;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
}

.forward-msg-preview-label {
  font-weight: 600;
  color: #666;
  margin-bottom: 4px;
}

.forward-msg-preview-content {
  color: #222;
  word-break: break-word;
}

.photo-indicator {
  font-style: italic;
  color: #555;
}

.modal-input {
  height: 44px;
  font-size: 14px;
  padding: 8px 14px;
  border-radius: 12px;
  border: 2px solid #555;
  background-color: #f5f5f5;
  box-sizing: border-box;
  width: 100%;
  transition: 0.25s ease;
}

.modal-input:focus {
  border-color: #000;
  background-color: #fff;
  outline: none;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.15);
}

.forward-chats-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 240px;
  overflow-y: auto;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 6px;
}

.forward-chat-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  border-radius: 10px;
  transition: background-color 0.15s ease;
}

.forward-chat-row:hover {
  background-color: #f5f5f5;
}

.forward-chat-avatar {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.forward-chat-meta {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.forward-chat-title {
  font-weight: 600;
  font-size: 14px;
  color: #222;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-header-badge {
  font-size: 11px;
  background-color: #eee;
  padding: 2px 6px;
  border-radius: 6px;
  color: #555;
}

.btn-action-small {
  padding: 6px 14px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: 0.15s ease;
}

.btn-primary-small {
  background-color: #222;
  color: white;
}

.btn-primary-small:hover:not(:disabled) {
  background-color: #444;
}

.btn-primary-small:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.no-chats-found {
  padding: 20px;
  text-align: center;
  color: #888;
  font-size: 14px;
}

.success-msg {
  color: #2e7d32;
  font-size: 14px;
  margin: 0;
  text-align: center;
  font-weight: 500;
}

.error-msg {
  color: red;
  font-size: 14px;
  margin: 0;
  text-align: center;
}

.modal-buttons {
  display: flex;
  justify-content: flex-end;
}

.btn-secondary {
  padding: 10px 16px;
  background-color: #ccc;
  color: #222;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-secondary:hover {
  background-color: #bbb;
}

@media (max-width: 600px) {
  .modal-forward-card {
    padding: 20px 16px;
    border-radius: 16px;
  }
}
</style>
