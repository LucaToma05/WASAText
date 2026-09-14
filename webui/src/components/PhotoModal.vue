<template>
  <!-- Se cliccato rimuove il pannello -->
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-photo-card">
      <div class="modal-card-header">
        <h2 class="modal-title">Invia foto</h2>
        <button type="button" class="btn-close-modal" @click="$emit('close')" title="Chiudi">✕</button>
      </div>

      <p class="modal-photo-desc">Incolla il link (URL) dell'immagine da inviare nella conversazione:</p>

      <!-- Anteprima della foto se inserita -->
      <div v-if="photoUrl.trim()" class="modal-photo-preview-wrap">
        <img 
          :src="photoUrl.trim()" 
          class="modal-photo-preview-img" 
          alt="Anteprima foto" 
          @error="onPreviewError"
        />
      </div>

      <!-- Inserisci URL -->
      <input 
        v-model="photoUrl"
        type="url"
        placeholder="https://esempio.com/foto.jpg"
        class="modal-input"
        @keyup.enter="handleSend"
      />

      <!-- Eventuali errori -->
      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

      <!-- Bottoni -->
      <div class="modal-buttons">
        <button 
          type="button" 
          class="btn-primary" 
          :disabled="isSending || !photoUrl.trim()" 
          @click="handleSend"
        >
          {{ isSending ? 'Invio in corso...' : 'Invia' }}
        </button>
        <button type="button" class="btn-secondary" @click="$emit('close')">Annulla</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "PhotoModal",
  props: {
    isSending: {
      type: Boolean,
      default: false
    }
  },
  emits: ["close", "send"],
  data() {
    return {
      photoUrl: "",
      errorMsg: ""
    };
  },
  methods: {
    // Mostra messaggio errore anteprima
    onPreviewError() {
      this.errorMsg = "Impossibile caricare l'anteprima dell'immagine. Verifica l'URL.";
    },
    // Gestisce invio foto
    handleSend() {
      const url = this.photoUrl.trim();
      if (!url) {
        this.errorMsg = "Inserisci l'URL di un'immagine";
        return;
      }
      this.errorMsg = "";
      this.$emit("send", url);
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

.modal-photo-card {
  background-color: #ffffff;
  padding: 30px 35px;
  border-radius: 22px;
  width: 460px;
  max-width: 95%;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.18);
  display: flex;
  flex-direction: column;
  gap: 18px;
  box-sizing: border-box;
}

.modal-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  font-size: 26px;
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

.modal-photo-desc {
  font-size: 14px;
  color: #555;
  margin: 0;
}

.modal-photo-preview-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  border-radius: 14px;
  padding: 10px;
  max-height: 200px;
  overflow: hidden;
}

.modal-photo-preview-img {
  max-width: 100%;
  max-height: 180px;
  object-fit: contain;
  border-radius: 8px;
}

.modal-input {
  width: 100%;
  height: 50px;
  font-size: 16px;
  padding: 10px 14px;
  border-radius: 12px;
  border: 2px solid #555;
  background-color: #f5f5f5;
  box-sizing: border-box;
  transition: 0.25s ease;
}

.modal-input:focus {
  border-color: #000;
  background-color: #fff;
  outline: none;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.15);
}

.error-msg {
  color: red;
  font-size: 14px;
  text-align: center;
  margin: 0;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.btn-primary {
  padding: 12px 20px;
  background-color: #222;
  color: white;
  border: none;
  border-radius: 14px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-primary:hover:not(:disabled) {
  background-color: #444;
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  padding: 12px 20px;
  background-color: #ccc;
  color: #222;
  border: none;
  border-radius: 14px;
  font-size: 16px;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-secondary:hover {
  background-color: #bbb;
}

@media (max-width: 600px) {
  .modal-photo-card {
    padding: 24px 20px;
    border-radius: 16px;
  }

  .modal-buttons {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
