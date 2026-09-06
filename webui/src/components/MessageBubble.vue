<template>
  <div class="message" :class="{ mine, 'is-deleted': isDeleted }">
    <span v-if="isDeleted" class="deleted-text"><em>Questo messaggio è stato eliminato</em></span>
    <div v-else-if="isPhoto" class="message-photo-wrapper">
      <img :src="photoUrl" class="message-photo-img" alt="Foto messaggio" loading="lazy" />
    </div>
    <span v-else>{{ textContent }}</span>

    <button 
      v-if="mine && !isDeleted" 
      type="button"
      class="btn-delete" 
      title="Elimina messaggio" 
      @click.stop.prevent="deleteMessageHandler"
    >
      🗑️
    </button>
  </div>
</template>

<script>
export default {
  props: {
    msg: Object,
    mine: Boolean,
    chatId: String
  },
  computed: {
    textContent() {
      if (!this.msg) return "";
      return this.msg.text || this.msg.content || "";
    },
    isDeleted() {
      if (!this.msg) return false;
      // Fallback: controllo legacy sull'isDeleted impostato localmente
      if (this.msg.isDeleted) return true;
      // Il backend imposta il testo al tombstone quando deleted=1
      const text = (this.textContent || "").trim();
      return text === "Questo messaggio \u00e8 stato eliminato";
    },
    isPhoto() {
      if (this.isDeleted) return false;
      if (this.msg && (this.msg.type === "photo" || this.msg.isPhoto)) {
        return true;
      }
      const text = this.textContent;
      if (!text || typeof text !== "string") return false;
      const t = text.trim();
      return (
        /^(https?:\/\/.*\.(?:png|jpg|jpeg|gif|webp|svg|bmp|avif)(\?.*)?)$/i.test(t) ||
        /^data:image\//i.test(t)
      );
    },
    photoUrl() {
      if (this.msg && this.msg.url) return this.msg.url;
      return this.textContent;
    }
  },
  methods: {
    async deleteMessageHandler() {
      console.log("[MessageBubble] Click elimina messaggio:", this.msg);
      this.$emit("delete", this.msg);

      if (this.chatId && this.msg?.id && this.$axios) {
        try {
          const token = localStorage.getItem("token");
          console.log(`[MessageBubble] Invio DELETE /chats/${this.chatId}/messages/${this.msg.id}`);
          await this.$axios.delete(`/chats/${this.chatId}/messages/${this.msg.id}`, {
            headers: { Authorization: "Bearer " + token }
          });
          console.log("[MessageBubble] Eliminazione completata con successo");
          // Soft-delete reattivo: aggiorna il testo per aggiornare subito la UI
          if (this.msg) {
            this.msg.text = "Questo messaggio \u00e8 stato eliminato";
            this.msg.isDeleted = true;
          }
        } catch (err) {
          console.error("[MessageBubble] Errore eliminazione messaggio:", err);
        }
      }
    }
  }
};
</script>

<style scoped>
.message {
  position: relative;
  max-width: 60%;
  padding: 10px;
  margin-bottom: 10px;
  background-color: #e6e6e6;
  border-radius: 8px;
  word-break: break-word;
}

.message.mine {
  background-color: #cce5ff;
  margin-left: auto;
}

.message.is-deleted {
  background-color: #f1f5f9;
  color: #94a3b8;
  border: 1px dashed #cbd5e1;
}

.deleted-text {
  font-style: italic;
  color: #94a3b8;
}

.message-photo-wrapper {
  display: flex;
  align-items: center;
}

.message-photo-img {
  display: block;
  max-width: 300px;
  width: 100%;
  height: auto;
  border-radius: 8px;
  object-fit: cover;
}

.btn-delete {
  background: none;
  border: none;
  cursor: pointer;
  padding: 2px 4px;
  margin-left: 6px;
  font-size: 14px;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.btn-delete:hover {
  opacity: 1;
}
</style>
