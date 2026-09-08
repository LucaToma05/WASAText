<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-card modal-group-card">
      <h2 class="modal-title">Nuovo Gruppo</h2>

      <!-- Nome Gruppo -->
      <div class="modal-form-group">
        <label class="modal-label">Nome del gruppo *</label>
        <input 
          v-model="groupName"
          type="text"
          placeholder="Nome del gruppo (3-25 caratteri)"
          class="modal-input modal-input-compact"
          maxlength="25"
          autofocus
        />
      </div>

      <!-- Foto Gruppo -->
      <div class="modal-form-group">
        <label class="modal-label">Foto del gruppo (opzionale)</label>
        <div class="photo-input-section">
          <div v-if="groupPhoto" class="photo-preview-box">
            <img :src="groupPhoto" class="group-photo-preview" alt="Anteprima foto" />
            <button type="button" class="btn-clear-photo" @click="groupPhoto = ''" title="Rimuovi foto">✕</button>
          </div>
          <div class="photo-inputs-wrapper">
            <input 
              v-model="groupPhoto"
              type="url"
              placeholder="Incolla il link (URL) dell'immagine..."
              class="modal-input modal-input-compact"
            />
          </div>
        </div>
      </div>

      <!-- Selezione Utenti -->
      <div class="modal-form-group">
        <label class="modal-label">Partecipanti</label>
        <div class="add-user-row">
          <input 
            v-model="userToAdd"
            type="text"
            placeholder="Username partecipante..."
            class="modal-input modal-input-compact"
            @keyup.enter.prevent="addParticipant"
          />
          <button type="button" class="btn-add-chip" @click="addParticipant">Aggiungi</button>
        </div>

        <!-- Membri selezionati -->
        <div class="selected-users-box">
          <span class="selected-title">Membri selezionati ({{ selectedParticipants.length }}):</span>
          <div v-if="selectedParticipants.length === 0" class="no-users-hint">
            Nessun altro utente selezionato. Tu verrai incluso automaticamente.
          </div>
          <div v-else class="chips-container">
            <span 
              v-for="user in selectedParticipants" 
              :key="user" 
              class="user-chip"
            >
              {{ user }}
              <button type="button" class="chip-remove-btn" @click="removeParticipant(user)" title="Rimuovi">✕</button>
            </span>
          </div>
        </div>
      </div>

      <p v-if="groupErrorMsg" class="error-msg">{{ groupErrorMsg }}</p>

      <div class="modal-buttons">
        <button class="btn-primary" :disabled="isCreatingGroup" @click="createGroup">
          {{ isCreatingGroup ? 'Creazione in corso...' : 'Crea gruppo' }}
        </button>
        <button type="button" class="btn-secondary" @click="$emit('close')">Annulla</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CreateGroupModal",
  emits: ["close", "created"],
  data() {
    return {
      groupName: "",
      groupPhoto: "",
      userToAdd: "",
      selectedParticipants: [],
      groupErrorMsg: "",
      isCreatingGroup: false
    };
  },
  methods: {
    addParticipant() {
      const u = this.userToAdd.trim();
      if (!u) return;
      if (u.length < 3 || u.length > 16) {
        this.groupErrorMsg = "Lo username deve essere tra 3 e 16 caratteri";
        return;
      }
      const myId = localStorage.getItem("userId") || localStorage.getItem("token");
      const myUsername = localStorage.getItem("username");
      if (u === myId || (myUsername && u.toLowerCase() === myUsername.toLowerCase())) {
        this.groupErrorMsg = "Tu sei già incluso automaticamente nel gruppo";
        return;
      }
      if (this.selectedParticipants.includes(u)) {
        this.groupErrorMsg = `L'utente "${u}" è già stato aggiunto`;
        return;
      }
      this.selectedParticipants.push(u);
      this.userToAdd = "";
      this.groupErrorMsg = "";
    },
    removeParticipant(username) {
      this.selectedParticipants = this.selectedParticipants.filter(u => u !== username);
    },
    async createGroup() {
      const name = this.groupName.trim();
      if (!name) {
        this.groupErrorMsg = "Inserisci il nome del gruppo";
        return;
      }
      if (name.length < 3 || name.length > 25) {
        this.groupErrorMsg = "Il nome del gruppo deve essere tra 3 e 25 caratteri";
        return;
      }

      this.isCreatingGroup = true;
      this.groupErrorMsg = "";

      try {
        const response = await this.$axios.post(
          "/chats",
          {
            type: true,
            isGroup: true,
            name: name,
            photo: this.groupPhoto.trim(),
            participants: this.selectedParticipants
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
        this.groupErrorMsg = err.response?.data || "Errore durante la creazione del gruppo";
      } finally {
        this.isCreatingGroup = false;
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

.modal-group-card {
  background-color: #ffffff;
  padding: 30px 35px;
  border-radius: 22px;
  width: 500px;
  max-width: 95%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  gap: 18px;
  box-sizing: border-box;
}

.modal-title {
  font-size: 28px;
  font-weight: bold;
  color: #222;
  text-align: center;
  margin: 0;
}

.modal-form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-align: left;
}

.modal-label {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.modal-input {
  width: 100%;
  box-sizing: border-box;
  border-radius: 12px;
  border: 2px solid #555;
  background-color: #f5f5f5;
  transition: 0.25s ease;
}

.modal-input:focus {
  border-color: #000;
  background-color: #fff;
  outline: none;
  box-shadow: 0 0 10px rgba(0,0,0,0.15);
}

.modal-input-compact {
  height: 46px;
  font-size: 15px;
  padding: 8px 14px;
}

.photo-input-section {
  display: flex;
  gap: 14px;
  align-items: center;
}

.photo-preview-box {
  position: relative;
  width: 50px;
  height: 50px;
  flex-shrink: 0;
}

.group-photo-preview {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #222;
}

.btn-clear-photo {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 18px;
  height: 18px;
  background: #ff4d4f;
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.photo-inputs-wrapper {
  flex: 1;
}

.add-user-row {
  display: flex;
  gap: 8px;
}

.btn-add-chip {
  padding: 0 16px;
  background-color: #222;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  white-space: nowrap;
  transition: 0.2s;
}

.btn-add-chip:hover {
  background-color: #444;
}

.selected-users-box {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 4px;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 12px;
  border: 1px dashed #ccc;
}

.selected-title {
  font-size: 12px;
  font-weight: 600;
  color: #555;
}

.no-users-hint {
  font-size: 12px;
  color: #888;
  font-style: italic;
}

.chips-container {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-height: 90px;
  overflow-y: auto;
}

.user-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background-color: #e6f4ff;
  color: #0050b3;
  border: 1px solid #91caff;
  padding: 4px 10px;
  border-radius: 14px;
  font-size: 13px;
  font-weight: 500;
}

.chip-remove-btn {
  background: none;
  border: none;
  color: #0050b3;
  font-size: 13px;
  font-weight: bold;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.chip-remove-btn:hover {
  color: #ff4d4f;
}

.error-msg {
  color: red;
  font-size: 15px;
  text-align: center;
  margin: 0;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-top: 6px;
}

.btn-primary {
  padding: 12px 18px;
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
  font-size: 16px;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-secondary:hover {
  background-color: #bbb;
}
</style>
