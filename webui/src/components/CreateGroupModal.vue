<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-card modal-group-card">

      <h2 class="modal-title">Nuovo Gruppo</h2>

      <!-- Nome del gruppo -->
      <div class="modal-form-group">
        <label class="modal-label">Nome del gruppo *</label>
        <input 
          v-model="groupName"
          type="text"
          placeholder="Nome del gruppo (3-25 caratteri)"
          class="modal-input modal-input-compact"
          maxlength="25"
        />
      </div>

      <!-- Foto del gruppo — solo URL -->
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

      <!-- Selezione Partecipanti con validazione real-time -->
      <div class="modal-form-group">
        <label class="modal-label">Partecipanti</label>
        <div class="add-user-row">
          <div class="user-input-wrapper">
            <input 
              v-model="userToAdd"
              type="text"
              placeholder="Username partecipante..."
              class="modal-input modal-input-compact"
              :class="{
                'input-valid':   validationState === 'valid',
                'input-invalid': validationState === 'invalid'
              }"
              @input="onUserInput"
              @keyup.enter.prevent="addParticipant"
            />
            <!-- Icona di stato inline sull'input -->
            <span v-if="validationState === 'checking'" class="input-icon input-icon--checking">⏳</span>
            <span v-else-if="validationState === 'valid'"   class="input-icon input-icon--valid">✓</span>
            <span v-else-if="validationState === 'invalid'" class="input-icon input-icon--invalid">✕</span>
          </div>
          <button
            type="button"
            class="btn-add-chip"
            :disabled="validationState !== 'valid'"
            @click="addParticipant"
          >
            Aggiungi
          </button>
        </div>

        <!-- Messaggio di feedback sotto l'input utente -->
        <transition name="fade">
          <p v-if="userSearchMsg" class="user-search-msg" :class="userSearchMsgClass">
            {{ userSearchMsg }}
          </p>
        </transition>

        <!-- Membri selezionati -->
        <div class="selected-users-box">
          <span class="selected-title">Membri selezionati ({{ selectedParticipants.length }}):</span>
          <div v-if="selectedParticipants.length === 0" class="no-users-hint">
            Nessun partecipante aggiunto. L'utente creatore sarà incluso automaticamente.
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

      <!-- Errore globale del form (nome gruppo, errori API) -->
      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

      <div class="modal-buttons">
        <button class="btn-primary" :disabled="loading" @click="createGroup">
          {{ loading ? 'Creazione in corso...' : 'Crea' }}
        </button>
        <button class="btn-secondary" @click="close">Annulla</button>
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

      // --- stato ricerca utente ---
      userToAdd: "",
      // 'idle' | 'checking' | 'valid' | 'invalid'
      validationState: "idle",
      userSearchMsg: "",
      debounceTimer: null,

      // --- partecipanti confermati ---
      selectedParticipants: [],

      // --- stato form ---
      errorMsg: "",
      loading: false
    };
  },

  computed: {
    userSearchMsgClass() {
      return {
        "msg-valid":   this.validationState === "valid",
        "msg-invalid": this.validationState === "invalid"
      };
    }
  },

  methods: {
    // ── Debounce sull'input username ────────────────────────────────
    onUserInput() {
      // Reset immediato ogni volta che l'utente digita
      clearTimeout(this.debounceTimer);
      this.validationState = "idle";
      this.userSearchMsg   = "";

      const u = this.userToAdd.trim();

      // Non fare nulla se l'input è troppo corto
      if (u.length < 3) {
        return;
      }

      // Avvia il debounce a 400ms
      this.validationState = "checking";
      this.userSearchMsg   = "Verifica in corso...";

      this.debounceTimer = setTimeout(() => {
        this.validateUsername(u);
      }, 400);
    },

    // ── Chiamata API per verificare l'esistenza dell'utente ─────────
    async validateUsername(username) {
      try {
        await this.$axios.get("/users", {
          params: { username },
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        // 200 → utente trovato
        if (this.selectedParticipants.includes(username)) {
          this.validationState = "invalid";
          this.userSearchMsg   = `"${username}" è già nella lista.`;
        } else {
          this.validationState = "valid";
          this.userSearchMsg   = `✓ Utente "${username}" trovato.`;
        }
      } catch (err) {
        if (err.response?.status === 404) {
          this.validationState = "invalid";
          this.userSearchMsg   = `Utente "${username}" non trovato.`;
        } else {
          this.validationState = "invalid";
          this.userSearchMsg   = "Errore durante la verifica. Riprova.";
        }
      }
    },

    // ── Aggiunge il partecipante solo se la validazione è positiva ──
    addParticipant() {
      if (this.validationState !== "valid") return;

      const u = this.userToAdd.trim();
      this.selectedParticipants.push(u);

      // Reset input e stato
      this.userToAdd       = "";
      this.validationState = "idle";
      this.userSearchMsg   = "";
      clearTimeout(this.debounceTimer);
    },

    removeParticipant(username) {
      this.selectedParticipants = this.selectedParticipants.filter(u => u !== username);
    },

    // ── Crea il gruppo ───────────────────────────────────────────────
    async createGroup() {
      if (this.groupName.trim() === "") {
        this.errorMsg = "Inserisci il nome del gruppo";
        return;
      }
      if (this.groupName.trim().length < 3 || this.groupName.trim().length > 25) {
        this.errorMsg = "Il nome del gruppo deve essere tra 3 e 25 caratteri";
        return;
      }

      this.loading  = true;
      this.errorMsg = "";

      try {
        const response = await this.$axios.post(
          "/chats",
          {
            type: true,
            isGroup: true,
            name: this.groupName.trim(),
            photo: this.groupPhoto.trim(),
            participants: this.selectedParticipants
          },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.$emit("created", response.data);
        this.close();
      } catch (err) {
        this.errorMsg = err.response?.data || "Errore durante la creazione del gruppo";
      } finally {
        this.loading = false;
      }
    },

    close() {
      clearTimeout(this.debounceTimer);
      this.groupName         = "";
      this.groupPhoto        = "";
      this.userToAdd         = "";
      this.validationState   = "idle";
      this.userSearchMsg     = "";
      this.selectedParticipants = [];
      this.errorMsg          = "";
      this.$emit("close");
    }
  },

  beforeUnmount() {
    // Pulizia timer al destroy del componente
    clearTimeout(this.debounceTimer);
  }
};
</script>

<style scoped>
/* ── Overlay ──────────────────────────────────────── */
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.45);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 500;
}

/* ── Card ─────────────────────────────────────────── */
.modal-card {
  background-color: #ffffff;
  padding: 30px;
  border-radius: 22px;
  width: 480px;
  max-width: 95%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.18);
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.modal-title {
  font-size: 28px;
  font-weight: bold;
  color: #222;
  text-align: center;
  margin: 0;
}

/* ── Form group ───────────────────────────────────── */
.modal-form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-align: left;
}

.modal-label {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

/* ── Input base ───────────────────────────────────── */
.modal-input {
  width: 100%;
  padding: 10px 14px;
  border-radius: 12px;
  border: 2px solid #555;
  background-color: #f5f5f5;
  font-size: 15px;
  box-sizing: border-box;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.modal-input-compact {
  height: 46px;
}

.modal-input:focus {
  border-color: #000;
  background-color: #fff;
  outline: none;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.15);
}

/* Stato valido/non valido sull'input */
.modal-input.input-valid {
  border-color: #22c55e;
  background-color: #f0fdf4;
}

.modal-input.input-invalid {
  border-color: #ef4444;
  background-color: #fff5f5;
}

/* ── Wrapper input utente con icona di stato ──────── */
.user-input-wrapper {
  position: relative;
  flex: 1;
}

.user-input-wrapper .modal-input {
  padding-right: 38px; /* spazio per l'icona */
}

.input-icon {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 15px;
  pointer-events: none;
}

.input-icon--checking { color: #f59e0b; }
.input-icon--valid    { color: #22c55e; }
.input-icon--invalid  { color: #ef4444; }

/* ── Messaggio feedback ricerca utente ────────────── */
.user-search-msg {
  font-size: 12px;
  margin: 0;
  padding: 0 4px;
  font-style: italic;
  color: #888;
}

.msg-valid   { color: #16a34a; }
.msg-invalid { color: #dc2626; }

/* ── Transizione fade ─────────────────────────────── */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ── Foto gruppo ──────────────────────────────────── */
.photo-input-section {
  display: flex;
  gap: 14px;
  align-items: center;
}

.photo-preview-box {
  position: relative;
  width: 60px;
  height: 60px;
  flex-shrink: 0;
}

.group-photo-preview {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #222;
}

.btn-clear-photo {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 20px;
  height: 20px;
  background: #ff4d4f;
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 11px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.photo-inputs-wrapper {
  flex: 1;
}

/* ── Riga aggiunta utente ─────────────────────────── */
.add-user-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.btn-add-chip {
  padding: 0 16px;
  height: 46px;
  background-color: #222;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  white-space: nowrap;
  transition: background-color 0.2s, opacity 0.2s;
  flex-shrink: 0;
}

.btn-add-chip:hover:not(:disabled) {
  background-color: #444;
}

.btn-add-chip:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* ── Partecipanti selezionati ─────────────────────── */
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

/* ── Errore globale ───────────────────────────────── */
.error-msg {
  color: #dc2626;
  font-size: 14px;
  text-align: center;
  margin: 0;
  font-weight: 500;
}

/* ── Bottoni azioni ───────────────────────────────── */
.modal-buttons {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.btn-primary {
  padding: 12px 24px;
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
