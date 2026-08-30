<template>
  <div class="overlay">
    <div class="overlay-box">
      <h3>Crea nuova chat</h3>

      <input 
        type="text"
        v-model="username"
        placeholder="Username utente"
      />

      <div class="overlay-actions">
        <button @click="createChat">Crea</button>
        <button @click="$emit('close')">Annulla</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: ""
    };
  },

  methods: {
    async createChat() {
      if (!this.username.trim()) return;

      try {
        await this.$axios.post("/createChat", {
          username: this.username
        });

        this.$emit("created");
        this.$emit("close");
      } catch (err) {
        console.log("Errore creazione chat:", err);
      }
    }
  }
};
</script>

<style scoped>
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.4);
  display: flex;
  justify-content: center;
  align-items: center;
}

.overlay-box {
  background: white;
  padding: 25px;
  border-radius: 8px;
  width: 350px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.overlay-actions {
  display: flex;
  justify-content: space-between;
}
</style>
