<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-card">

      <h2 class="modal-title">Crea nuova chat</h2>

      <input 
        v-model="username"
        type="text"
        placeholder="Username dell'utente"
        class="modal-input"
      />

      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

      <div class="modal-buttons">
        <button class="btn-primary" @click="createChat">Crea chat</button>
        <button class="btn-secondary" @click="close">Annulla</button>
      </div>

    </div>
  </div>
</template>

<script>
export default {
  emits: ["close", "created"],

  data() {
    return {
      username: "",
      errorMsg: ""
    };
  },

  methods: {
    async createChat() {
      if (this.username.trim() === "") {
        this.errorMsg = "Inserisci un username";
        return;
      }

      try {
        const userId = localStorage.getItem("userId");

        const response = await this.$axios.post(`/chats`, {
          creatorId: userId,
          targetUsername: this.username
        });

        this.$emit("created", response.data);
        this.close();

      } catch (err) {
        this.errorMsg = err.response?.data?.error || "Errore creazione chat";
      }
    },

    close() {
      this.$emit("close");
    }
  }
};
</script>

<style scoped>
/* Overlay */
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

/* Card */
.modal-card {
  background-color: #ffffff;
  padding: 40px;
  border-radius: 22px;
  width: 420px;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.modal-title {
  font-size: 32px;
  font-weight: bold;
  color: #222;
  text-align: center;
}

.modal-input {
  width: 100%;
  height: 60px;
  font-size: 22px;
  padding: 12px;
  border-radius: 14px;
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

.error-msg {
  color: red;
  font-size: 18px;
  text-align: center;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
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

.btn-primary:hover {
  background-color: #444;
  transform: translateY(-3px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.25);
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
</style>
