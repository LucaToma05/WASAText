<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-card">
      <h2 class="modal-title">Crea nuova chat</h2>

      <input 
        v-model="username"
        type="text"
        placeholder="Username dell'utente"
        class="modal-input"
        autofocus
        @keyup.enter="createChat"
      />

      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

      <div class="modal-buttons">
        <button class="btn-primary" :disabled="loading" @click="createChat">
          {{ loading ? 'Creazione...' : 'Crea chat' }}
        </button>
        <button class="btn-secondary" type="button" @click="$emit('close')">Annulla</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CreateChatModal",
  emits: ["close", "created"],
  data() {
    return {
      username: "",
      errorMsg: "",
      loading: false
    };
  },
  methods: {
    async createChat() {
      const target = this.username.trim();
      if (!target) {
        this.errorMsg = "Inserisci un username";
        return;
      }

      this.loading = true;
      this.errorMsg = "";

      try {
        const response = await this.$axios.post(
          "/chats",
          {
            type: false,
            name: target
          },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        const newChatId = response.data.chatId || response.data.id;
        this.$emit("created", newChatId);
        this.$emit("close");
      } catch (err) {
        this.errorMsg = err.response?.data || "Errore creazione chat";
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.45);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 500;
}

.modal-card {
  background-color: #ffffff;
  padding: clamp(24px, 5vw, 40px);
  border-radius: 22px;
  width: 90%;
  max-width: 420px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  gap: 25px;
  box-sizing: border-box;
}

.modal-title {
  font-size: clamp(22px, 5vw, 32px);
  font-weight: bold;
  color: #222;
  text-align: center;
}

.modal-input {
  width: 100%;
  height: clamp(48px, 8vw, 60px);
  font-size: clamp(16px, 3vw, 20px);
  padding: 12px;
  border-radius: 14px;
  border: 2px solid #555;
  background-color: #f5f5f5;
  box-sizing: border-box;
  transition: 0.25s ease;
}

.modal-input:focus {
  border-color: #000;
  background-color: #fff;
  outline: none;
  box-shadow: 0 0 10px rgba(0,0,0,0.15);
}

.error-msg {
  color: red;
  font-size: 16px;
  text-align: center;
  margin: 0;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.btn-primary {
  padding: 12px 18px;
  background-color: #222;
  color: white;
  border: none;
  border-radius: 14px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-primary:hover:not(:disabled) {
  background-color: #444;
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.2);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  padding: 12px 18px;
  background-color: #ccc;
  border: none;
  border-radius: 14px;
  font-size: 18px;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-secondary:hover {
  background-color: #bbb;
}

@media (max-width: 600px) {
  .modal-card {
    border-radius: 16px;
    gap: 18px;
  }

  .modal-buttons {
    flex-direction: column;
    gap: 10px;
  }

  .modal-buttons .btn-primary,
  .modal-buttons .btn-secondary {
    width: 100%;
    text-align: center;
  }
}
</style>
