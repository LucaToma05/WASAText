<template>
  <!-- Se cliccata rimuove il pannello -->
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-group-card">
      <div class="modal-header-row">
        <h2 class="modal-title">Gestione Gruppo</h2>
        <button type="button" class="btn-close-modal" @click="$emit('close')" title="Chiudi">✕</button>
      </div>

      <!-- Modifica Nome Gruppo -->
      <div class="modal-form-group">
        <label class="modal-label">Nome del gruppo</label>
        <div class="input-with-action">
          <input 
            v-model="editGroupName"
            type="text"
            placeholder="Nome del gruppo (3-25 caratteri)"
            class="modal-input"
            maxlength="25"
          />
          <button 
            type="button" 
            class="btn-action-small"
            :disabled="isSavingGroupName || !editGroupName.trim() || editGroupName.trim() === chat?.name"
            @click="updateGroupName"
          >
            {{ isSavingGroupName ? 'Salvataggio...' : 'Aggiorna' }}
          </button>
        </div>
      </div>

      <!-- Modifica Foto Gruppo (solo URL) -->
      <div class="modal-form-group">
        <label class="modal-label">Foto del gruppo (URL immagine)</label>
        <div class="photo-edit-section">
          <!-- Preview della foto -->
          <div v-if="editGroupPhoto" class="photo-preview-box">
            <img :src="editGroupPhoto" class="group-photo-preview" alt="Anteprima foto" />
            <button type="button" class="btn-clear-photo" @click="editGroupPhoto = ''" title="Rimuovi foto">✕</button>
          </div>
          <div class="input-with-action">
            <input 
              v-model="editGroupPhoto"
              type="url"
              placeholder="Incolla il link (URL) dell'immagine..."
              class="modal-input"
            />
            <button 
              type="button" 
              class="btn-action-small"
              :disabled="isSavingGroupPhoto || !editGroupPhoto.trim() || editGroupPhoto.trim() === chat?.photo"
              @click="updateGroupPhoto"
            >
              {{ isSavingGroupPhoto ? 'Salvataggio...' : 'Aggiorna' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Aggiungi Nuovi Partecipanti -->
      <div class="modal-form-group">
        <label class="modal-label">Aggiungi nuovo membro (cerca utente WASAText)</label>
        <div class="input-with-action">
          <input 
            v-model="newMemberUsername"
            type="text"
            placeholder="Username partecipante..."
            class="modal-input"
            @keyup.enter="addMemberToGroup"
          />
          <button 
            type="button" 
            class="btn-action-small btn-primary-small"
            :disabled="isAddingMember || !newMemberUsername.trim()"
            @click="addMemberToGroup"
          >
            {{ isAddingMember ? 'Aggiunta...' : 'Aggiungi' }}
          </button>
        </div>
      </div>

      <!-- Feedback stato -->
      <p v-if="successMsg" class="success-msg">{{ successMsg }}</p>
      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

      <hr class="modal-divider" />

      <!-- Abbandona gruppo -->
      <div class="modal-danger-zone">
        <button 
          type="button" 
          class="btn-danger"
          :disabled="isLeavingGroup"
          @click="leaveGroup"
        >
          {{ isLeavingGroup ? 'Uscita in corso...' : 'Abbandona Gruppo' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "GroupSettingsModal",
  props: {
    chat: {
      type: Object,
      required: true
    },
    currentUserId: {
      type: String,
      default: ""
    }
  },
  emits: ["close", "nameUpdated", "photoUpdated", "memberAdded", "leftGroup"],
  data() {
    return {
      editGroupName: this.chat?.name || "",
      editGroupPhoto: this.chat?.photo || "",
      newMemberUsername: "",
      errorMsg: "",
      successMsg: "",
      isSavingGroupName: false,
      isSavingGroupPhoto: false,
      isAddingMember: false,
      isLeavingGroup: false
    };
  },
  methods: {
    async updateGroupName() {
      // Controllo del nome del gruppo
      const name = this.editGroupName.trim();
      if (!name) {
        this.errorMsg = "Inserisci un nome per il gruppo";
        return;
      }
      if (name.length < 3 || name.length > 25) {
        this.errorMsg = "Il nome del gruppo deve essere tra 3 e 25 caratteri";
        return;
      }

      this.isSavingGroupName = true;
      this.errorMsg = "";
      this.successMsg = "";

      // Aggiorna nome del gruppo
      try {
        await this.$axios.patch(
          `/chats/${this.chat.chatId}/name`,
          { name },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.successMsg = "Nome del gruppo aggiornato con successo!";
        this.$emit("nameUpdated", name);
      } catch (err) {
        this.errorMsg = err.response?.data || "Errore durante l'aggiornamento del nome";
      } finally {
        this.isSavingGroupName = false;
      }
    },

    async updateGroupPhoto() {
      // Controllo URL foto
      const photoUrl = this.editGroupPhoto.trim();
      if (!photoUrl) {
        this.errorMsg = "Inserisci un URL valido per la foto del gruppo";
        return;
      }
      if (photoUrl.length < 3) {
        this.errorMsg = "La foto specificata non è valida";
        return;
      }

      this.isSavingGroupPhoto = true;
      this.errorMsg = "";
      this.successMsg = "";

      // Aggiorna foto
      try {
        await this.$axios.put(
          `/chats/${this.chat.chatId}/photo`,
          { url: photoUrl },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.successMsg = "Foto del gruppo aggiornata con successo!";
        this.$emit("photoUpdated", photoUrl);
      } catch (err) {
        this.errorMsg = err.response?.data || "Errore durante l'aggiornamento della foto";
      } finally {
        this.isSavingGroupPhoto = false;
      }
    },

    async addMemberToGroup() {
      // Controllo username
      const username = this.newMemberUsername.trim();
      if (!username) {
        this.errorMsg = "Inserisci un username";
        return;
      }
      if (username.length < 3 || username.length > 16) {
        this.errorMsg = "Lo username deve essere tra 3 e 16 caratteri";
        return;
      }

      // Controllo se l'utente è già nel gruppo
      const myId = this.currentUserId || localStorage.getItem("token");
      const myUsername = localStorage.getItem("username");
      if (username === myId || (myUsername && username.toLowerCase() === myUsername.toLowerCase())) {
        this.errorMsg = "Fai già parte di questo gruppo!";
        return;
      }

      this.isAddingMember = true;
      this.errorMsg = "";
      this.successMsg = "";

      // Aggiungi membro
      try {
        await this.$axios.post(
          `/chats/${this.chat.chatId}/members`,
          { userId: username },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.successMsg = `Utente "${username}" aggiunto al gruppo con successo!`;
        this.newMemberUsername = "";
        this.$emit("memberAdded", username);
      } catch (err) {
        this.errorMsg = err.response?.data || `Errore nell'aggiunta dell'utente "${username}"`;
      } finally {
        this.isAddingMember = false;
      }
    },

    async leaveGroup() {
      this.isLeavingGroup = true;
      this.errorMsg = "";

      // Abbandona il gruppo
      try {
        await this.$axios.delete(
          `/chats/${this.chat.chatId}/members`,
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.$emit("leftGroup", this.chat.chatId);
        this.$emit("close");
      } catch (err) {
        this.errorMsg = err.response?.data || "Errore durante l'uscita dal gruppo";
      } finally {
        this.isLeavingGroup = false;
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

.modal-group-card {
  background-color: #ffffff;
  padding: 30px 35px;
  border-radius: 22px;
  width: 500px;
  max-width: 95%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.18);
  display: flex;
  flex-direction: column;
  gap: 18px;
  box-sizing: border-box;
}

.modal-header-row {
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
  height: 44px;
  font-size: 15px;
  padding: 8px 14px;
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
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.15);
}

.input-with-action {
  display: flex;
  gap: 8px;
}

.input-with-action .modal-input {
  flex: 1;
}

.btn-action-small {
  padding: 0 16px;
  background-color: #eee;
  color: #222;
  border: 1px solid #ccc;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: 0.15s ease;
}

.btn-action-small:hover:not(:disabled) {
  background-color: #ddd;
}

.btn-action-small:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary-small {
  background-color: #222;
  color: white;
  border-color: #222;
}

.btn-primary-small:hover:not(:disabled) {
  background-color: #444;
}

.photo-edit-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.photo-preview-box {
  position: relative;
  width: 50px;
  height: 50px;
  flex-shrink: 0;
}

.group-photo-preview {
  width: 100%;
  height: 100%;
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

.modal-divider {
  border: 0;
  height: 1px;
  background-color: #eee;
  margin: 4px 0;
}

.modal-danger-zone {
  display: flex;
  justify-content: flex-end;
}

.btn-danger {
  padding: 10px 18px;
  background-color: #ffebee;
  color: #c62828;
  border: 1px solid #ffcdd2;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s ease;
}

.btn-danger:hover:not(:disabled) {
  background-color: #c62828;
  color: white;
}

.btn-danger:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 600px) {
  .modal-group-card {
    padding: 24px 20px;
    border-radius: 16px;
  }
}
</style>
